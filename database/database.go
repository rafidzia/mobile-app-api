package database

import (
	"database/sql"
	"errors"
	"os"
)

func OpenDatabase() (*sql.DB, error) {
	var mysql_addr string
	if os.Getenv("MYSQL_HOST") != "" && os.Getenv("MYSQL_PORT") != "" && os.Getenv("MYSQL_DB") != "" && os.Getenv("MYSQL_USER") != "" {
		mysql_addr = os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DB")
	} else {
		return nil, errors.New("MYSQL_HOST, MYSQL_PORT, MYSQL_DB, MYSQL_USER must be set")
	}
	db, err := sql.Open("mysql", mysql_addr)

	return db, err
}
