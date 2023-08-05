package gosqldb

import "database/sql"

type SqlDatabase interface {
	// DbInit Performs initialization operation for the database.
	DbInit(driver, connectorString string, maxRetries int)

	// DbClose Performs cleanup operation for the database.
	DbClose()

	// RunQuery Runs thr QueryString to the database
	RunQuery(query string, args ...interface{}) (*sql.Rows, error)

	// RunTransaction Runs all the queries as a transaction
	RunTransaction(queries []Query) error

	BeginTransaction() (*sql.Tx, error)
}

type Query struct {
	QueryString string
	Args        []interface{}
}
