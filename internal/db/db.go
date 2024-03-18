package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DBExecutor struct {
	Db *sql.DB
}

// InitDB initializes the database with the given connection string.
func InitDB(dataSourceName string) (*DBExecutor, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping the database: %v", err)
	}

	log.Println("Successfully connected to the database")

	return &DBExecutor{Db: db}, nil
}
