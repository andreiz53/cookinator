package database

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/andreiz53/cookinator/util"
	"github.com/jackc/pgx/v5"
)

var testQueries *Queries
var testDB *pgx.Conn

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("could not load env:", err)
	}
	testDB, err = pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("could not connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
