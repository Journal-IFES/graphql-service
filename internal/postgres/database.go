package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitPostgresDB(connStr string) error {
	d, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	db = d

	return nil
}

func GetPostgresDB() *sql.DB {
	return db
}

func ClosePostgresDB() {
	db.Close()
}
