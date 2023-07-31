# Go SQL Database Module

This Go lang module simplifies running queries to SQL-based databases. It provides an interface and implementations for managing database connections and executing queries.

## Installation

To install this module using `go get`, you can run the following command:

```
go get -u github.com/hk1997/gosqldb
```

Or, you can import it in your Go project:

```go
import (
    "github.com/hk1997/gosqldb"
)
```

### Interface
The module provides the following interface:

```
import "database/sql"

type SqlDatabase interface {
    // DbInit performs initialization operation for the database.
    DbInit(driver, connectorString string, maxRetries int)

    // DbClose performs cleanup operation for the database.
    DbClose()

    // RunQuery runs the query on the database.
    RunQuery(query string, args ...interface{}) (*sql.Rows, error)
}

```

### Usage
The module includes two implementations of the SqlDatabase interface: `PostgresImpl` and `MockSqlDatabase`.


#### PostgresSQL Database
To initialize a PostgreSQL database, you can use the PostgresImpl implementation:

Create a global database accessor to be used across your module. This will help to share the same 
database object across multiple packages.
```
package appModule

import (
	db "github.com/hk1997/gosqldb"
)

var appDb db.SqlDatabase

func InitDatabase(database db.SqlDatabase) {
	appDb = database
}

func InitPostgresDatabase(driver, connectorString string, maxRetries int) {
	var boozingoDatabase db.SqlDatabase
	appDb = &db.PostgresImpl{}
	appDb.DbInit(driver, connectorString, maxRetries)
	InitDatabase(appDb)
}
```

In your main.go, initialize the postgres database:
```
package main

import (
	"api/appModule"
	"api/server"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	appModule.InitPostgresDatabase(appModule.DB_DRIVER, appModule.DB_CONNECTOR_STRING, appModule.DB_RETRIES)

	server.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("Server started on port 8080")
	log.Fatal(srv.ListenAndServe())
}
```


##### Running queries

To run the query simply import the appModule and call database accessor:
```
func SendOtp(extension, mobile string) (bool, error) {
	otp := generateOTP()
	retries := 3

	rows, err := appModule.appDb.RunQuery(INSERT_OTP_QUERY, extension, mobile, otp, retries)
	if err != nil {
		fmt.Println("Error running query", err)
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		var success bool
		err := rows.Scan(&success)
		if err != nil {
			return false, err
		}
		return success, nil
	}

	return false, nil
}
```

##### Testing

* Initialize mock database using sql mocking library before each test

```
import (
	"api/appModule"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	database "github.com/hk1997/gosqldb"
	"os"
	"testing"
)

var db *sql.DB
var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	db, mock, _ = sqlmock.New()
	defer db.Close()
	var mockDb database.SqlDatabase = &database.MockSqlDatabase{Db: db}
	appModule.InitDatabase(mockDb)

	exitCode := m.Run()

	os.Exit(exitCode)
}
```

* Write the test to mock the result of the query

```
func TestSendOtp(t *testing.T) {

	// Set up the expected query and result for SendOtp
	mockRows := sqlmock.NewRows([]string{"retriesZero"}).AddRow(true)
	mock.ExpectQuery(regexp.QuoteMeta(INSERT_OTP_QUERY)).WithArgs("test_extension", "test_mobile", sqlmock.AnyArg(), 3).
		WillReturnRows(mockRows)

	// Call the SendOtp function with test data
	retriesZero, err := SendOtp("test_extension", "test_mobile")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !retriesZero {
		t.Error("Expected retriesZero to be true")
	}

	// Check that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
```

### License
This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.