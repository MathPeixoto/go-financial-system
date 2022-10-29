package db

import (
	"database/sql"
	"github.com/MathPeixoto/go-financial-system/util"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../app.env")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	testDb, err = sql.Open(config.DatabaseDriver, config.DatabaseSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
