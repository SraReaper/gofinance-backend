package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:postgres@localhost:5432/go_finance?sslmode=disable"
)

var testQueries *Queries

// TestMain roda antes dos testes
func TestMain(m *testing.M) {
	conectando, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	testQueries = New(conectando)
	os.Exit(m.Run())
}
