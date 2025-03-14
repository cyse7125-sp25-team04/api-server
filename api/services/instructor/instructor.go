package instructor

import (
	"errors"
	"fmt"
	"time"
	"webapp/db"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Instructor represents the instructor schema.
type Instructor struct {
	InstructorId int       `gorm:"primaryKey;autoIncrement;column:instructor_id" json:"instructorId"`
	UserId       int       `gorm:"column:user_id" json:"userId"`
	DepartmentId int       `gorm:"column:department_id" json:"departmentId"`
	CreatedAt    time.Time `gorm:"column:date_created" json:"date_created"`
	UpdatedAt    time.Time `gorm:"column:date_modified" json:"date_modified"`
}

// CreateInstructor inserts a new instructor record into the Instructors table.
func CreateInstructor(inst *Instructor) error {
	// Get the MySQL connection.
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return errors.New("error connecting to MySQL")
	}
	// Create the new instructor record.
	if err := gormDB.Create(inst).Error; err != nil {
		fmt.Println(err)
		return errors.New("error inserting new instructor")
	}

	return nil
}

// UpdateInstructor performs a full update of an instructor record.
func UpdateInstructor(inst *Instructor) error {
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return errors.New("error connecting to MySQL")
	}
	// Update all fields of the instructor record.
	if err := gormDB.Model(&Instructor{}).Where("instructor_id = ?", inst.InstructorId).Updates(inst).Error; err != nil {
		return errors.New("error updating instructor")
	}
	return nil
}

// PatchInstructor performs a partial update on an instructor record.
func PatchInstructor(instructorId int, fields map[string]interface{}) error {
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
	if err := gormDB.Model(&Instructor{}).Where("instructor_id = ?", instructorId).Updates(fields).Error; err != nil {
		return errors.New("error patching instructor")
	}
	return nil
}

// DeleteInstructor removes an instructor record by its ID.
func DeleteInstructor(instructorId int) error {
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
	if err := gormDB.Delete(&Instructor{}, instructorId).Error; err != nil {
		return errors.New("error deleting instructor")
	}
	return nil
}

func isInstructorExistsByUserId(userId int) (bool, error) {
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return false, errors.New("error connecting to MySQL")
	}
	var inst Instructor
	result := gormDB.Where("user_id = ?", userId).First(&inst)
	if result.Error != nil {
		return false, nil
	}
	return true, nil
}

func GetInstructorById(instructorId int) (*Instructor, error) {
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return nil, errors.New("error connecting to MySQL")
	}
	var inst Instructor
	result := gormDB.Where("instructor_id = ?", instructorId).First(&inst)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("error getting instructor")
	}
	return &inst, nil
}
