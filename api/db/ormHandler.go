package db

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetOrmDatabase() (*gorm.DB, error) {
	if gormDBInstance != nil {
		return gormDBInstance, nil
	}

	sqlDB, err := GetMySQLConn()
	if err != nil {
		log.WithError(err).Error("Failed to connect to MySQL")
		return nil, errors.New("error connecting to MySQL")
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.WithError(err).Error("Failed to initialize GORM")
		return nil, errors.New("error initializing GORM")
	}

	log.Info("Successfully initialized GORM database connection")
	gormDBInstance = gormDB
	return gormDBInstance, nil
}

func CloseDB() {
	if sqlDBInstance != nil {
		sqlDBInstance.Close()
		log.Info("MySQL connection closed")
	}
}
