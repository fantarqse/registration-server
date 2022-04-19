package db

import "database/sql"

func New() (*sql.DB, error) {
	return sql.Open(
		"postgres",
		"user=postgres password=qwerty dbname=registration_db sslmode=disable",
	)
}
