package database

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/tonievictor/dotenv"
)

func TestDB(t *testing.T) {
	dotenv.Config("../../.env")

	t.Run("Test correct database creation", func(t *testing.T) {
		_, err := Init("postgres", os.Getenv("DB_CONN_STRING"))
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Test incorrect database creation", func(t *testing.T) {
		_, err := Init("postgr", os.Getenv("DB_CONN_STRI"))
		if err == nil {
			t.Error(err)
		}
	})
}
