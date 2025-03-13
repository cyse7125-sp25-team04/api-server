package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	course "webapp/services/courses"
	healthcheck "webapp/services/healthcheck"
	instructor "webapp/services/instructor"
	trace "webapp/services/trace"
	user "webapp/services/user"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Hello, World!")
	startServer()
}

func startServer() {
	fmt.Println("Server started...")

	// Create a new router.
	router := mux.NewRouter()

	// Register endpoints with appropriate HTTP methods.
	router.HandleFunc("/healthz", healthcheckHandler).Methods("GET")
	router.HandleFunc("/v1/user", userHandler).Methods("POST")
	// GET endpoint to retrieve user details (using basic auth).
	router.HandleFunc("/v1/user", getUserHandler).Methods("GET")
	// PUT endpoint for updating a user (only allowed after auth).
	router.HandleFunc("/v1/user", updateUserHandler).Methods("PUT")
	// New endpoint for creating an instructor; only accessible to admin.
	router.HandleFunc("/v1/instructor", createInstructorHandler).Methods("POST")
	router.HandleFunc("/v1/instructor", updateInstructorHandler).Methods("PUT")
	router.HandleFunc("/v1/instructor", patchInstructorHandler).Methods("PATCH")
	router.HandleFunc("/v1/instructor", deleteInstructorHandler).Methods("DELETE")
	// New endpoints for /v1/course.
	router.HandleFunc("/v1/course", createCourseHandler).Methods("POST")
	router.HandleFunc("/v1/course", deleteCourseHandler).Methods("DELETE")
	// router.HandleFunc("/v1/course/{course_id}/", trace.TraceHandler).Methods("GET")
	// router.HandleFunc("/v1/instructor/{instructor_id}/", trace.TraceHandler).Methods("GET")
	//Trace Handling
	router.HandleFunc("/v1/course/{course_id}/trace", trace.TraceHandler).Methods("POST")
	// router.HandleFunc("/v1/course/{course_id}/trace", trace.TraceHandler).Methods("GET")
	// router.HandleFunc("/v1/course/{course_id}/trace/{trace_id}", trace.TraceHandler).Methods("GET")
	// router.HandleFunc("/v1/course/{course_id}/trace/{trace_id}", trace.TraceHandler).Methods("DELETE")
	// End points for departments and terms.
	// add department

	//add term
	// Start the server on port 8080.
	http.ListenAndServe(":8080", router)
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.ContentLength > 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := healthcheck.Check()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	var newUser user.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received user: %+v\n", newUser)
	if newUser.FirstName == "" || newUser.LastName == "" || newUser.Username == "" || newUser.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	err = user.CreateUser(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// getUserHandler authenticates via basic auth and returns user details.
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	authUser, err := user.Authenticate(username, password)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authUser)
}

// updateUserHandler authenticates the user and then updates the user's data.
// Returns 400 Bad Request for invalid payload/path, 401 for unauthorized,
// and 204 No Content on success.
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Verify basic auth credentials.
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	authUser, err := user.Authenticate(username, password)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ensure the request has the expected Content-Type.
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	// Validate that the URL path is exactly "/v1/user".
	if r.URL.Path != "/v1/user" {
		http.Error(w, "Bad Request: invalid path", http.StatusBadRequest)
		return
	}

	// Decode the JSON payload into a User struct.
	var updatedUser user.User
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Force the update to use the authenticated user's ID.
	updatedUser.ID = authUser.ID

	// Update the user record.
	err = user.UpdateUser(&updatedUser)
	if err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// On success, return 204 No Content.
	w.WriteHeader(http.StatusNoContent)
}

// createInstructorHandler is accessible only to admin (username "admin", pwd "admin").
// It decodes the JSON payload for a new instructor, appends the admin user's ID dynamically,
// and inserts it into the database.
func createInstructorHandler(w http.ResponseWriter, r *http.Request) {
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
	var inst instructor.Instructor
	fmt.Println("r.Body", r.Body)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&inst)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	fmt.Printf("inst", inst)
	// Create the instructor record.
	err = instructor.CreateInstructor(&inst)
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
func updateInstructorHandler(w http.ResponseWriter, r *http.Request) {
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
	// Decode JSON payload into an Instructor struct.
	var inst instructor.Instructor
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&inst)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// Require instructorId to update the record.
	if inst.InstructorId == 0 {
		http.Error(w, "Missing required field: instructorId", http.StatusBadRequest)
		return
	}
	// Perform full update.
	err = instructor.UpdateInstructor(&inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the updated instructor.
	w.WriteHeader(http.StatusNoContent)
}

// patchInstructorHandler handles partial update requests (PATCH) for instructors.
func patchInstructorHandler(w http.ResponseWriter, r *http.Request) {
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
	err = instructor.PatchInstructor(instID, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with no content.
	w.WriteHeader(http.StatusNoContent)
}

// deleteInstructorHandler handles DELETE requests for instructors.
func deleteInstructorHandler(w http.ResponseWriter, r *http.Request) {
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
	err = instructor.DeleteInstructor(payload.InstructorId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with no content.
	w.WriteHeader(http.StatusNoContent)
}

// createCourseHandler handles POST requests to create a new course.
func createCourseHandler(w http.ResponseWriter, r *http.Request) {
	// Require admin authentication.
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

	// Decode JSON payload.
	var crs course.Course
	fmt.Println(r.Body)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&crs); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// Validate required fields.
	if crs.Name == "" {
		http.Error(w, "Missing required fields: name and course_code", http.StatusBadRequest)
		return
	}

	// Create the course.
	if err := course.CreateCourse(&crs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 201 Created with the created course data.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(crs)
}

// deleteCourseHandler handles DELETE requests to remove a course based on its course code.
func deleteCourseHandler(w http.ResponseWriter, r *http.Request) {
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

	// Expect a JSON payload with course_code.
	var payload struct {
		CourseCode int `json:"course_code"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if payload.CourseCode == 0 {
		http.Error(w, "Missing required field: course_code", http.StatusBadRequest)
		return
	}

	// Retrieve the course by its course code.
	crs, err := course.GetCourseByCode(payload.CourseCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Delete the course using its CourseId.
	if err := course.DeleteCourse(crs.CourseId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content on successful deletion.
	w.WriteHeader(http.StatusNoContent)
}
