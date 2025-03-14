package course

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webapp/services/constants"
	"webapp/services/user"

	"github.com/gorilla/mux"
)

// createCourseHandler handles POST requests to create a new course.
func CreateCourseHandler(w http.ResponseWriter, r *http.Request) {
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
	var crs Course
	fmt.Println(r.Body)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&crs); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// Validate required fields.
	if crs.Name == "" || crs.CourseCode == "" || crs.DepartmentId == 0 {
		http.Error(w, "Missing required fields: name, courseCode and DepartmentId", http.StatusBadRequest)
		return
	}
	dbcrs, err := GetCourseByCode(crs.CourseCode)
	if err != nil {
		http.Error(w, "Course with given course code already exists", http.StatusInternalServerError)
		return
	}
	if dbcrs != nil {
		http.Error(w, "Course with given course code already exists", http.StatusConflict)
		return
	}
	// Create the course.
	if err := CreateCourse(&crs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 201 Created with the created course data.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(crs)
}

func UpdateCourseHandler(w http.ResponseWriter, r *http.Request) {

	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	adminUser, err := user.Authenticate(username, password)
	if err != nil || adminUser.Role != constants.AdminRole {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	courseIdStr := mux.Vars(r)["course_id"]
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}
	var crs Course
	fmt.Println(r.Body)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&crs); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// Validate required fields.
	if crs.Name == "" || crs.CourseCode == "" || crs.DepartmentId == 0 {
		http.Error(w, "Missing required fields: name, courseCode and DepartmentId", http.StatusBadRequest)
		return
	}
	dbcrs, err := GetCourseById(courseId)
	if dbcrs == nil || err != nil {
		http.Error(w, "Course not found", http.StatusBadRequest)
		return
	}
	if dbcrs.CourseCode != crs.CourseCode {
		http.Error(w, "can not update Course code", http.StatusBadRequest)
		return
	}
	err = UpdateCourse(&crs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// deleteCourseHandler handles DELETE requests to remove a course based on its course code.
func DeleteCourseHandler(w http.ResponseWriter, r *http.Request) {
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
		CourseCode string `json:"courseCode"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if payload.CourseCode == "" {
		http.Error(w, "Missing required field: courseCode", http.StatusBadRequest)
		return
	}

	// Retrieve the course by its course code.
	crs, err := GetCourseByCode(payload.CourseCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Delete the course using its CourseId.
	if err := DeleteCourse(crs.CourseId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content on successful deletion.
	w.WriteHeader(http.StatusNoContent)
}

func GetCourseHandler(w http.ResponseWriter, r *http.Request) {
	courseIdStr := mux.Vars(r)["course_id"]
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}
	crs, err := GetCourseById(courseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(crs)
}
