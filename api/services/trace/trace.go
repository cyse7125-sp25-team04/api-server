package trace

import (
	"errors"
	"fmt"
	"mime/multipart"
	"time"
	"webapp/config"
	"webapp/db"
	gcpgateway "webapp/services/gcpGateway"

	"gorm.io/gorm"
)

type Report struct {
	ID              int       `gorm:"primaryKey;autoIncrement;column:report_id" json:"traceId"`
	Type            string    `gorm:"column:type" json:"type"`
	EnrollmentCount int       `gorm:"column:enrollment_count" json:"enrollmentCount"`
	ResponsesCount  int       `gorm:"column:responses_count" json:"responsesCount"`
	InstructorId    int       `gorm:"column:instructor_id" json:"instructorId"`
	FileName        string    `gorm:"column:file_name" json:"filename"`
	CourseId        int       `gorm:"column:course_id" json:"courseId"`
	TermId          int       `gorm:"column:term_id" json:"termId"`
	BucketPath      string    `gorm:"column:bucket_path" json:"bucketPath"`
	DateCreated     time.Time `gorm:"column:date_created" json:"dateCreated"`
}

func addTraceReport(folderPath string, file multipart.File, fileName string) error {
	bucketName := config.GetEnv("STORAGE_BUCKET_NAME")
	if err := gcpgateway.UploadFile(bucketName, folderPath, fileName, file); err != nil {
		return err
	}
	return nil
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

func GetTracesByCourseId(courseId int) ([]Report, error) {
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return nil, errors.New("error connecting to database")
	}
	var traces []Report
	if err := gormDB.Where("course_id = ?", courseId).Find(&traces).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, errors.New("error fetching traces")
	}
	return traces, nil
}

func GetTraceByCourseIdAndTraceId(courseId int, traceId int) (*Report, error) {
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return nil, errors.New("error connecting to database")
	}
	var trace Report
	if err := gormDB.Where("course_id = ? AND report_id = ?", courseId, traceId).First(&trace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("error fetching trace")
	}
	return &trace, nil
}

func DeleteTraceByCourseIdAndTraceId(courseId int, traceId int) error {
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return errors.New("error connecting to database")
	}
	if err := gormDB.Where("course_id = ? AND report_id = ?", courseId, traceId).Delete(&Report{}).Error; err != nil {
		return errors.New("error deleting trace")
	}
	return nil
}
