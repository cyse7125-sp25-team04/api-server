package trace

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"webapp/config"
	"webapp/db"
	constants "webapp/services/constants"
	user "webapp/services/user"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	option "google.golang.org/api/option"
)

type Report struct {
	ID              int       `grom:"primaryKey;autoIncrement;column:report_id" json:"traceId"`
	Type            string    `grom:"column:type" json:"type"`
	EnrollmentCount int       `grom:"column:enrollment_count" json:"enrollmentCount"`
	ResponsesCount  int       `grom:"column:responses_count" json:"responsesCount"`
	InstructorId    int       `grom:"column:instructor_id" json:"instructorId"`
	FileName        string    `grom:"column:file_name" json:"filename"`
	CourseId        int       `grom:"column:course_id" json:"courseId"`
	TermId          int       `grom:"column:term_id" json:"termId"`
	BucketPath      string    `grom:"column:bucket_path" json:"bucketPath"`
	DateCreated     time.Time `grom:"column:date_created" json:"dateCreated"`
}

func TraceHandler(w http.ResponseWriter, r *http.Request) {
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
	// 3.forward the request
	addTraceReport(w, r)
}

func addTraceReport(w http.ResponseWriter, r *http.Request) {
	// 1.get the trace report
	file, handler, err := r.FormFile("traceFile")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("filename: %v\n", file)
	// 2. Get additional form data
	termIDStr := r.FormValue("termId")
	instructorIDStr := r.FormValue("instructorId")
	reportType := r.FormValue("reportType")
	vars := mux.Vars(r)
	courseIdStr := vars["course_id"]
	fmt.Fprintf(w, "Course ID: %s", courseIdStr)
	fmt.Println("term Id: ", termIDStr)
	fmt.Println("instructorIDStr: ", instructorIDStr)
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

	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		http.Error(w, "Invalid course Id", http.StatusBadRequest)
		return
	}
	// 2. get gcp context
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithQuotaProject(config.GetEnvConfig().GOOGLE_PROJECT_ID))
	if err != nil {
		http.Error(w, "Failed to create Cloud Storage client", http.StatusInternalServerError)
		return
	}

	defer client.Close()
	// 3. get the bucker objects
	bucketName := config.GetEnv("STORAGE_BUCKET_NAME")
	bucket := client.Bucket(bucketName)
	folderPath := fmt.Sprintf("/%s/%d/%d/%s/", courseIdStr, instructorID, termID, reportType)
	object := bucket.Object(folderPath + handler.Filename)
	writer := object.NewWriter(ctx)

	// 4. Copy the file contents from the HTTP request to GCS
	if _, err := io.Copy(writer, file); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to upload file to Cloud Storage", http.StatusInternalServerError)
		return
	}
	if err := writer.Close(); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to finalize the upload", http.StatusInternalServerError)
		return
	}
	//5. add the details to the database.
	if err := addEntrytoDatabase(termID, instructorID, courseId, handler.Filename, folderPath, reportType); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to add entry to database", http.StatusInternalServerError)
		return
	}
	// 6. Success response
	fmt.Fprintf(w, "File uploaded successfully to bucket: %s", bucketName)
}

func addEntrytoDatabase(termID int, instructorID int, courseId int, filename, bucketPath string, reportType string) error {
	// 1. get database connection

	grom, err := db.GetOrmDatabase()
	if err != nil {
		return errors.New("error setting up connection to database")
	}
	// 2. create a new report entry
	report := Report{
		Type:         reportType,
		TermId:       termID,
		InstructorId: instructorID,
		CourseId:     courseId,
		FileName:     filename,
		BucketPath:   bucketPath,
		DateCreated:  time.Now().UTC(),
	}
	fmt.Println("report: ", report)
	if err := grom.Create(&report).Error; err != nil {
		return errors.New("error inserting new report")
	}
	return nil
}
