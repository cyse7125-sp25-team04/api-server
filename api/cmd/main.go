package main

import (
	"fmt"
	"net/http"

	course "webapp/services/courses"
	"webapp/services/healthcheck"
	"webapp/services/instructor"
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
	router.HandleFunc("/healthz", healthcheck.HealthcheckHandler).Methods("GET")
	router.HandleFunc("/v1/user", user.UserHandler).Methods("POST")
	// GET endpoint to retrieve user details (using basic auth).
	router.HandleFunc("/v1/user", user.GetUserHandler).Methods("GET")
	// PUT endpoint for updating a user (only allowed after auth).
	router.HandleFunc("/v1/user", user.UpdateUserHandler).Methods("PUT")

	// New endpoint for creating an instructor; only accessible to admin.
	router.HandleFunc("/v1/instructor", instructor.CreateInstructorHandler).Methods("POST")
	router.HandleFunc("/v1/instructor/{instructor_id}", instructor.GetInstrutorHandler).Methods("GET")
	router.HandleFunc("/v1/instructor/{instructor_id}", instructor.UpdateInstructorHandler).Methods("PUT")
	router.HandleFunc("/v1/instructor", instructor.PatchInstructorHandler).Methods("PATCH")
	router.HandleFunc("/v1/instructor", instructor.DeleteInstructorHandler).Methods("DELETE")

	// New endpoints for /v1/course.
	router.HandleFunc("/v1/course", course.CreateCourseHandler).Methods("POST")
	router.HandleFunc("/v1/course", course.DeleteCourseHandler).Methods("DELETE")
	router.HandleFunc("/v1/course/{course_id}", course.UpdateCourseHandler).Methods("PUT")
	router.HandleFunc("/v1/course/{course_id}", course.GetCourseHandler).Methods("GET")

	//Trace Handling
	router.HandleFunc("/v1/course/{course_id}/trace", trace.UploadTraceHandler).Methods("POST")
	router.HandleFunc("/v1/course/{course_id}/trace", trace.GetAllTracesHandler).Methods("GET") // get all traces
	router.HandleFunc("/v1/course/{course_id}/trace/{trace_id}", trace.GetTraceHandler).Methods("GET")
	router.HandleFunc("/v1/course/{course_id}/trace/{trace_id}", trace.DeleteTraceHandler).Methods("DELETE")

	// End points for departments and terms.
	// add department

	//add term
	// Start the server on port 8080.
	http.ListenAndServe(":8080", router)
}
