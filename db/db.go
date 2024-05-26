package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

func Init() {
	var err error
	Database, err = sql.Open("mysql", "admin:admin@1234@tcp(13.232.54.105)/flutter_forge")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ensure the database connection is available
	if err := Database.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Optionally, set up connection pooling
	Database.SetMaxOpenConns(25)
	Database.SetMaxIdleConns(25)
	Database.SetConnMaxLifetime(5 * time.Minute)
}
