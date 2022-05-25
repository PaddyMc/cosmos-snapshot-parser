package database

import (
	"database/sql"

	_ "github.com/lib/pq" // nolint
)

func GetDBConnection() (*sql.DB, error) {
	//	sslMode := "disable"
	//	schema := "public"
	//	host := "localhost"
	//	port := "5432"
	//	dbname := "chain"
	//	user := "plural"
	//	password := "plural"
	maxOpen := 1
	maxIdle := 1

	connStr := "postgresql://plural:plural@localhost:5432/chain?sslmode=disable"

	postgresDb, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Set max open connections
	postgresDb.SetMaxOpenConns(maxOpen)
	postgresDb.SetMaxIdleConns(maxIdle)

	return postgresDb, nil
}
