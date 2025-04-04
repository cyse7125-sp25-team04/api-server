package db

import (
	"database/sql"
	"fmt"
	"sync"
	"webapp/config"
	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)
var(
	sqlDBInstance  *sql.DB
	gormDBInstance *gorm.DB
	once           sync.Once
)
func GetMySQLConn() (*sql.DB, error) {
	// Format the MySQL connection string
	once.Do(func() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.GetEnvConfig().DB_USERNAME,
		config.GetEnvConfig().DB_PASSWORD,
		config.GetEnvConfig().DB_HOST,
		config.GetEnvConfig().DB_PORT,
		config.GetEnvConfig().DB_NAME,
	)
	fmt.Println(dsn)
	// Open the connection to MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		db.Close()
		return
	}

	// Set the schema (if necessary)
	_, err = db.Exec(fmt.Sprintf("USE %s", config.GetEnvConfig().DB_NAME))
	if err != nil {
		fmt.Println(err)
		db.Close() // Close the connection if setting schema fails
		return
	}

		db.SetMaxOpenConns(100) // Adjust based on your MySQL server's max_connections
		db.SetMaxIdleConns(10)  // Keep some idle connections for reuse
		db.SetConnMaxLifetime(0)
		sqlDBInstance = db
	})
	if sqlDBInstance == nil {
		return nil, fmt.Errorf("failed to create MySQL connection")
	}
	return sqlDBInstance, nil
}
