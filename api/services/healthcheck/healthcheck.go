package healthcheck

import (
	"errors"
	"fmt"
	"time"
	"webapp/db"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// HealthCheckRecord represents a record in the webapp table.
type HealthCheckRecord struct {
	Check_id   uint      `gorm:"primaryKey"`
	Check_time time.Time `json:"datetime"`
}

// TableName sets the insert table name for this struct type.
func (HealthCheckRecord) TableName() string {
	return "Healthchecks"
}

// Check inserts a new healthcheck record using GORM.
func Check() error {
	// Get the MySQL connection from your db package.
	sqlDB, err := db.GetMySQLConn()
	if err != nil {
		fmt.Println(err)
		log.WithError(err).Error("error connecting to MySQL")
		return errors.New("error connecting to MySQL")
	}

	// Wrap the *sql.DB connection with GORM using the MySQL driver.
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		log.WithError(err).Error("error initializing GORM")
		return errors.New("error initializing GORM")
	}

	// Create a new healthcheck record with the current UTC time.
	record := HealthCheckRecord{
		Check_time: time.Now().UTC(),
	}
	if err := gormDB.Create(&record).Error; err != nil {
		fmt.Println(err)
		log.WithError(err).Error("error inserting into healthchecks table")
		return errors.New("error inserting into healthchecks table")
	}

	return nil
}
