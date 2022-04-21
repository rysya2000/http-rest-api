package sqlstore_test

import (
	"os"
	"testing"
)

var databaseURL string

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// databaseURL = "host=localhost port= dbname=postgres sslmode=disable"
		databaseURL = "host=localhost port=5432 user=postgres password=Qwerty123 dbname=postgres sslmode=disable"
	}

	os.Exit(m.Run())
}
