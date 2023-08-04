package gosqldb

import (
	"database/sql"
	"fmt"
	"log"
)

type MockSqlDatabase struct {
	Db *sql.DB
}

func (mdb *MockSqlDatabase) DbInit(_, _ string, _ int) {

}

func (mdb *MockSqlDatabase) DbClose() {
	mdb.Db.Close()
}

func (mdb *MockSqlDatabase) RunQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := mdb.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (mdb *MockSqlDatabase) RunTransaction(queries []Query) error {
	log.Println("Running transaction")
	tx, err := mdb.Db.Begin()
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
