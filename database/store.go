package database

import (
	"database/sql"

	_ "github.com/lib/pq" // nolint
)

func GetDBConnection(
	connStr string,
) (*sql.DB, error) {
	maxOpen := 3
	maxIdle := 3

	postgresDb, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Set max open connections
	postgresDb.SetMaxOpenConns(maxOpen)
	postgresDb.SetMaxIdleConns(maxIdle)

	return postgresDb, nil
}
