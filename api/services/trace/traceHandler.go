package trace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webapp/services/constants"
	"webapp/services/user"

	"github.com/gorilla/mux"
)

func UploadTraceHandler(w http.ResponseWriter, r *http.Request) {
	// 1. validate the request.
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "File too large or invalid request", http.StatusBadRequest)
		return
	}
	// 2. authenticate the user
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	authUser, err := user.Authenticate(username, password)
	if err != nil || authUser.Role != constants.AdminRole {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 3. Get additional form data
	file, handler, err := r.FormFile("traceFile")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	reportType := r.FormValue("reportType")
	termIDStr := r.FormValue("termId")
	instructorIDStr := r.FormValue("instructorId")
	vars := mux.Vars(r)
	courseIdStr := vars["course_id"]
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		http.Error(w, "Invalid courseId", http.StatusBadRequest)
		return
	}
	termID, err := strconv.Atoi(termIDStr)
	if err != nil {
		http.Error(w, "Invalid termId", http.StatusBadRequest)
		return
	}

	instructorID, err := strconv.Atoi(instructorIDStr)
	if err != nil {
		http.Error(w, "Invalid instructorId", http.StatusBadRequest)
		return
	}

	folderPath := fmt.Sprintf("/%s/%d/%d/%s/", courseIdStr, instructorID, termID, reportType)
	// 3.forward the request
	if err := addTraceReport(folderPath, file, handler.Filename); err != nil {
		http.Error(w, "Failed to upload file. Error :"+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := addEntrytoDatabase(termID, instructorID, courseId, handler.Filename, folderPath, reportType); err != nil {
		fmt.Println(err)
		// TODO: Delete the Uploaded file
		http.Error(w, "Failed to add entry to database", http.StatusInternalServerError)
		return
	}
}

// URL: /v1/course/{course_id}/trace
func GetAllTracesHandler(w http.ResponseWriter, r *http.Request) {
	// 1. get curse id
	courseIdStr := mux.Vars(r)["course_id"]
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		http.Error(w, "Invalid courseId", http.StatusBadRequest)
		return
	}
	// 2. get traces by course id

	traces, err := GetTracesByCourseId(courseId)
	// 3. return traces
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(traces)
}

// URL:/v1/course/{course_id}/trace/{trace_id}
func GetTraceHandler(w http.ResponseWriter, r *http.Request) {
	// 1. get course id and trace id
	courseIdStr := mux.Vars(r)["course_id"]
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		http.Error(w, "Invalid courseId", http.StatusBadRequest)
		return
	}
	traceIdStr := mux.Vars(r)["trace_id"]
	traceId, err := strconv.Atoi(traceIdStr)
	if err != nil {
		http.Error(w, "Invalid traceId", http.StatusBadRequest)
		return
	}
	// 2. get trace by course id and trace id
	trace, err := GetTraceByCourseIdAndTraceId(courseId, traceId)
	// 3. return trace
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(trace)

}

// URL: /v1/course/{course_id}/trace/{trace_id}
func DeleteTraceHandler(w http.ResponseWriter, r *http.Request) {
	// 1. get course id and trace id
	courseIdStr := mux.Vars(r)["course_id"]
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		http.Error(w, "Invalid courseId", http.StatusBadRequest)
		return
	}
	traceIdStr := mux.Vars(r)["trace_id"]
	traceId, err := strconv.Atoi(traceIdStr)
	if err != nil {
		http.Error(w, "Invalid traceId", http.StatusBadRequest)
		return
	}
	// 2. delete trace by course id and trace id
	err = DeleteTraceByCourseIdAndTraceId(courseId, traceId)
	// 3. return response
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
