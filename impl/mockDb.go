package impl

import "database/sql"

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
