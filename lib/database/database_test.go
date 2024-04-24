package database

import (
	"bytes"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/tonie-ng/go-dotenv"
)

func TestDB(t *testing.T) {
	dotenv.Config("../../.env")
	var bw bytes.Buffer
	l := log.New(&bw, "nest_test: ", log.LstdFlags)

	t.Run("Test correct database creation", func(t *testing.T) {
		db := New("postgres", os.Getenv("DB_CONN_STRING"), l)
		err := db.Init()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Test incorrect database creation", func(t *testing.T) {
		db := New("post", os.Getenv("DB_CONN_STRING"), l)
		err := db.Init()
		if err == nil {
			t.Error(err)
		}
	})
}
