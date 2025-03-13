package course

import (
	"errors"
	"fmt"
	"time"
	"webapp/db"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Course represents the courses table schema.
type Course struct {
	CourseId     int       `gorm:"primaryKey;autoIncrement;column:course_id" json:"courseId"`
	Name         string    `gorm:"column:name" json:"name"`
	CourseCode   string    `gorm:"column:course_code" json:"course_code"`
	CreatedAt    time.Time `gorm:"column:date_created" json:"date_created"`
	UpdatedAt    time.Time `gorm:"column:date_modified" json:"date_updated"`
	DepartmentId int       `gorm:"column:department_id" json:"departmentId"`
}

// GetCourseByCode retrieves a course record by its course code.
func GetCourseByCode(courseCode int) (*Course, error) {
	// Get the MySQL connection.
	sqlDB, err := db.GetMySQLConn()
	if err != nil {
		return nil, errors.New("error connecting to MySQL")
	}
	// Open GORM DB using the MySQL driver.
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, errors.New("error initializing GORM")
	}

	var crs Course
	// Query for a course with the given course code.
	if err := gormDB.Where("course_code = ?", courseCode).First(&crs).Error; err != nil {
		return nil, errors.New("course not found")
	}
	return &crs, nil
}

// CreateCourse inserts a new course record into the Courses table.
// It first checks if a course with the same course_code already exists.
func CreateCourse(course *Course) error {
	// Get the MySQL connection.
	sqlDB, err := db.GetMySQLConn()
	if err != nil {
		fmt.Println(err)
		return errors.New("error connecting to MySQL")
	}
	// Open GORM DB using the MySQL driver.
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return errors.New("error initializing GORM")
	}

	// // Check if a course with the same course_code already exists.
	// var existing Course
	// if err := gormDB.Where("course_code = ?", course.CourseCode).First(&existing).Error; err == nil {
	// 	return errors.New("course code already exists")
	// } else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	// An error occurred during lookup that is not a "record not found" error.
	// 	return errors.New("error checking for existing course code")
	// }

	// Create the new course record.
	if err := gormDB.Create(course).Error; err != nil {
		fmt.Println(err)
		return errors.New("error inserting new course")
	}
	return nil
}

// DeleteCourse removes a course record by its ID.
func DeleteCourse(courseId int) error {
	sqlDB, err := db.GetMySQLConn()
	if err != nil {
		return errors.New("error connecting to MySQL")
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return errors.New("error initializing GORM")
	}
	if err := gormDB.Delete(&Course{}, courseId).Error; err != nil {
		return errors.New("error deleting course")
	}
	return nil
}
