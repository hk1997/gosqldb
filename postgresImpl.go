package gosqldb

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type PostgresImpl struct {
	db *sql.DB
}

func (postgres *PostgresImpl) DbInit(driver string, connectorString string, maxRetries int) {
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		postgres.db, err = sql.Open(driver, connectorString)
		if err == nil {
			log.Println("Connected to the database successfully")
			return // Connection successful, exit the loop
		}

		log.Printf("Attempt %d: Failed to connect to the DB: %v", attempt, err)

		// Calculate the backoff duration using exponential backoff strategy
		backoffDuration := time.Duration(1<<uint(attempt)) * time.Second

		// Wait before retrying the connection
		log.Printf("Attempt %d: Retrying in %v...", attempt, backoffDuration)
		time.Sleep(backoffDuration)
	}

	log.Fatalf("Failed to connect to the database after %d attempts", maxRetries)
}

func (postgres *PostgresImpl) DbClose() {
	if postgres.db != nil {
		log.Println("Closing Db connection")
		postgres.db.Close()
	}
}

func (postgres *PostgresImpl) RunQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := postgres.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (postgres *PostgresImpl) RunTransaction(queries []Query) error {
	log.Println("Running transaction")
	tx, err := postgres.db.Begin()
	if err != nil {
		return err
	}

	commitErr := func() error {
		if err != nil {
			log.Println("Transaction failed with error:", err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("transaction rollback error: %v\noriginal error: %v", rollbackErr, err)
			}
			return err
		}
		if commitErr := tx.Commit(); commitErr != nil {
			return fmt.Errorf("transaction commit error: %v", commitErr)
		}
		return nil
	}

	for _, q := range queries {
		_, err = tx.Exec(q.QueryString, q.Args...)
		if err != nil {
			return commitErr()
		}
	}

	return commitErr()
}

func (postgres *PostgresImpl) BeginTransaction() (*sql.Tx, error) {
	return postgres.db.Begin()
}
