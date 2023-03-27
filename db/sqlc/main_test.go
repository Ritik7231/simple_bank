package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbSource = "postgresql://root:ritik@localhost:5432/simple_bank?sslmode=disable"
	dbDriver = "postgres"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can't connect to db: ", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
