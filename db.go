package gosqldb

import "database/sql"

type SqlDatabase interface {
	// DbInit Performs initialization operation for the database.
	DbInit(driver, connectorString string, maxRetries int)

	// DbClose Performs cleanup operation for the database.
	DbClose()

	// RunQuery Runs thr query to the database
	RunQuery(query string, args ...interface{}) (*sql.Rows, error)
}
