package instructor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webapp/services/constants"
	"webapp/services/user"

	"github.com/gorilla/mux"
)

// createInstructorHandler is accessible only to admin (username "admin", pwd "admin").
// It decodes the JSON payload for a new instructor, appends the admin user's ID dynamically,
// and inserts it into the database.
func CreateInstructorHandler(w http.ResponseWriter, r *http.Request) {
	// Check for basic auth credentials.
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Authenticate the admin user.
	adminUser, err := user.Authenticate(username, password)
	if err != nil || adminUser.Username != "admin" {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate Content-Type.
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	// Decode the JSON payload.
	var inst Instructor
	fmt.Println("r.Body", r.Body)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&inst)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	fmt.Printf("inst", inst)
	// Create the instructor record.
	if u, err := user.GetUserByID(strconv.Itoa(inst.UserId)); err != nil || u.Role != constants.InstructorRole {
		http.Error(w, "user not found or user is not an instructor", http.StatusBadRequest)
	}

	exists, err := isInstructorExistsByUserId(inst.UserId)
	if err != nil {
		http.Error(w, "Error checking instructor existence: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Instructor already exists", http.StatusBadRequest)
		return
	}
	err = CreateInstructor(&inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 201 Created with the created instructor data.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inst)
}

// updateInstructorHandler handles full update requests (PUT) for instructors.
func UpdateInstructorHandler(w http.ResponseWriter, r *http.Request) {
	// Admin authentication.
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	instUser, err := user.Authenticate(username, password)
	if err != nil || (!(instUser.Username == "admin" || instUser.Username == username)) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Validate Content-Type.
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	// Decode JSON payload into an Instructor struct.
	var inst Instructor
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&inst)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if inst.UserId != 0 && instUser.ID != inst.UserId {
		http.Error(w, "Unauthorized to update User Id", http.StatusBadRequest)
		return
	}
	// Require instructorId to update the record.
	muxVars := mux.Vars(r)
	instIDStr := muxVars["instructor_id"]
	instID, err := strconv.Atoi(instIDStr)
	if err != nil {
		http.Error(w, "Invalid instructor ID", http.StatusBadRequest)
		return
	}
	inst.InstructorId = instID
	// Perform full update.
	err = UpdateInstructor(&inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the updated instructor.
	w.WriteHeader(http.StatusNoContent)
}

// patchInstructorHandler handles partial update requests (PATCH) for instructors.
func PatchInstructorHandler(w http.ResponseWriter, r *http.Request) {
	// Admin authentication.
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	adminUser, err := user.Authenticate(username, password)
	if err != nil || adminUser.Username != "admin" {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Validate Content-Type.
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	// Decode JSON payload into a generic map.
	var payload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// Ensure instructorId is provided.
	idVal, exists := payload["instructorId"]
	if !exists {
		http.Error(w, "Missing required field: instructorId", http.StatusBadRequest)
		return
	}
	// JSON numbers are decoded as float64.
	idFloat, ok := idVal.(float64)
	if !ok {
		http.Error(w, "Invalid type for instructorId", http.StatusBadRequest)
		return
	}
	instID := int(idFloat)
	// Remove instructorId from payload to avoid updating the primary key.
	delete(payload, "instructorId")
	if len(payload) == 0 {
		http.Error(w, "No fields provided for update", http.StatusBadRequest)
		return
	}
	// Perform partial update.
	err = PatchInstructor(instID, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with no content.
	w.WriteHeader(http.StatusNoContent)
}

// deleteInstructorHandler handles DELETE requests for instructors.
func DeleteInstructorHandler(w http.ResponseWriter, r *http.Request) {
	// Admin authentication.
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	adminUser, err := user.Authenticate(username, password)
	if err != nil || adminUser.Username != "admin" {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Validate Content-Type.
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	// Expect a JSON payload with the instructorId to delete.
	var payload struct {
		InstructorId int `json:"instructorId"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if payload.InstructorId == 0 {
		http.Error(w, "Missing required field: instructorId", http.StatusBadRequest)
		return
	}
	// Perform deletion.
	err = DeleteInstructor(payload.InstructorId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with no content.
	w.WriteHeader(http.StatusNoContent)
}

func GetInstrutorHandler(w http.ResponseWriter, r *http.Request) {
	// Public endPoint
	// Expect a JSON payload with the instructorId to fetch.
	instructorIdStr := mux.Vars(r)["instructor_id"]
	instructorId, err := strconv.Atoi(instructorIdStr)
	if err != nil {
		http.Error(w, "Invalid instructorId", http.StatusBadRequest)
		return
	}
	// Perform fetch.
	inst, err := GetInstructorById(instructorId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if inst == nil {
		http.Error(w, "Instructor not found", http.StatusBadRequest)
		return
	}
	// Return the instructor.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inst)
}
