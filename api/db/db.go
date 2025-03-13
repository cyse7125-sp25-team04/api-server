package db

import (
	"database/sql"
	"fmt"

	"webapp/config"
	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

func GetMySQLConn() (*sql.DB, error) {
	// Format the MySQL connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true",
		config.GetEnvConfig().DB_USERNAME,
		config.GetEnvConfig().DB_PASSWORD,
		config.GetEnvConfig().DB_HOST,
		config.GetEnvConfig().DB_PORT)
	fmt.Println(dsn)
	// Open the connection to MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Set the schema (if necessary)
	_, err = db.Exec(fmt.Sprintf("USE %s", config.GetEnvConfig().DB_NAME))
	if err != nil {
		fmt.Println(err)
		db.Close() // Close the connection if setting schema fails
		return nil, err
	}

	return db, nil
}
