package pgstore_test

import (
	"iryzzh/practicum-gophermart/cmd/gophermart/config"
	"log"
	"os"
	"testing"
)

var (
	dsn string
)

func TestMain(m *testing.M) {
	cfg, err := config.New()
	if err != nil {
		log.Println("config failed:", err)
		os.Exit(1)
	}

	dsn = cfg.DatabaseURI
	if dsn == "" {
		dsn = "postgresql://postgres:postgres@localhost/praktikum?sslmode=disable"
	}

	os.Exit(m.Run())
}
