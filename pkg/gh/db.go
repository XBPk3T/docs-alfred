package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var respositorySql = `
CREATE TABLE IF NOT EXISTS repository (
    id varchar(255) PRIMARY KEY,
    url varchar(255) NOT NULL,
    user varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    description text,
    pushed_at timestamp,
    created_at timestamp,
    updated_at timestamp
)`

func OpenDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(respositorySql); err != nil {
		return nil, err
	}

	return db, nil
}
