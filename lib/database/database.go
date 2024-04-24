package database

import (
	"database/sql"
	"log"
)

type Database struct {
	db         *sql.DB
	driver     string
	datasource string
	logger     *log.Logger
}

func New(driver, connstring string, l *log.Logger) *Database {
	database := &Database{
		driver:     driver,
		datasource: connstring,
		logger:     l,
	}

	return database
}

func (database *Database) Init() error {
	db, err := sql.Open(database.driver, database.datasource)
	if err != nil {
		database.logger.Println(err)
		return err
	}

	database.db = db

	err = db.Ping()
	if err != nil {
		database.logger.Println(err)
		return err
	}

	return nil
}
