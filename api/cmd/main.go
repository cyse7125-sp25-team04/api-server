package main

import (
	"net/http"
	"os"

	course "webapp/services/courses"
	"webapp/services/healthcheck"
	"webapp/services/instructor"
	trace "webapp/services/trace"
	user "webapp/services/user"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Configure logrus for structured logging (JSON format).
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.Info("Application initialization starting...")
	startServer()
}

func startServer() {
	log.Info("Server setup initiated")

	// Create a new router.
	router := mux.NewRouter()

	// Register endpoints with appropriate HTTP methods.
	router.HandleFunc("/healthz", healthcheck.HealthcheckHandler).Methods("GET")
	log.Info("Registered endpoint: GET /healthz")

	// User routes
	router.HandleFunc("/v1/user", user.UserHandler).Methods("POST")
	log.Info("Registered endpoint: POST /v1/user")

	router.HandleFunc("/v1/user", user.GetUserHandler).Methods("GET")
	log.Info("Registered endpoint: GET /v1/user")

	router.HandleFunc("/v1/user", user.UpdateUserHandler).Methods("PUT")
	log.Info("Registered endpoint: PUT /v1/user")

	// Instructor routes
	router.HandleFunc("/v1/instructor", instructor.CreateInstructorHandler).Methods("POST")
	log.Info("Registered endpoint: POST /v1/instructor")

	router.HandleFunc("/v1/instructor/{instructor_id}", instructor.GetInstrutorHandler).Methods("GET")
	log.Info("Registered endpoint: GET /v1/instructor/{instructor_id}")

	router.HandleFunc("/v1/instructor/{instructor_id}", instructor.UpdateInstructorHandler).Methods("PUT")
	log.Info("Registered endpoint: PUT /v1/instructor/{instructor_id}")

	router.HandleFunc("/v1/instructor", instructor.PatchInstructorHandler).Methods("PATCH")
	log.Info("Registered endpoint: PATCH /v1/instructor")

	router.HandleFunc("/v1/instructor", instructor.DeleteInstructorHandler).Methods("DELETE")
	log.Info("Registered endpoint: DELETE /v1/instructor")

	// Course routes
	router.HandleFunc("/v1/course", course.CreateCourseHandler).Methods("POST")
	log.Info("Registered endpoint: POST /v1/course")

	router.HandleFunc("/v1/course", course.DeleteCourseHandler).Methods("DELETE")
	log.Info("Registered endpoint: DELETE /v1/course")

	router.HandleFunc("/v1/course/{course_id}", course.UpdateCourseHandler).Methods("PUT")
	log.Info("Registered endpoint: PUT /v1/course/{course_id}")

	router.HandleFunc("/v1/course/{course_id}", course.GetCourseHandler).Methods("GET")
	log.Info("Registered endpoint: GET /v1/course/{course_id}")

	// Trace Handling
	router.HandleFunc("/v1/course/{course_id}/trace", trace.UploadTraceHandler).Methods("POST")
	log.Info("Registered endpoint: POST /v1/course/{course_id}/trace")

	router.HandleFunc("/v1/course/{course_id}/trace", trace.GetAllTracesHandler).Methods("GET")
	log.Info("Registered endpoint: GET /v1/course/{course_id}/trace")

	router.HandleFunc("/v1/course/{course_id}/trace/{trace_id}", trace.GetTraceHandler).Methods("GET")
	log.Info("Registered endpoint: GET /v1/course/{course_id}/trace/{trace_id}")

	router.HandleFunc("/v1/course/{course_id}/trace/{trace_id}", trace.DeleteTraceHandler).Methods("DELETE")
	log.Info("Registered endpoint: DELETE /v1/course/{course_id}/trace/{trace_id}")

	// Start the server.
	log.WithField("port", "8080").Info("Server is listening")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.WithError(err).Fatal("Server failed to start")
	}
}
