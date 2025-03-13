package db

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetOrmDatabase() (*gorm.DB, error) {
	sqlDB, err := GetMySQLConn()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error connecting to MySQL")
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error initializing GORM")
	}

	return gormDB, nil
}
