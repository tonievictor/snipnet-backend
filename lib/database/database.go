package database

import (
	"database/sql"
)

func Init(driver, datasource string) (*sql.DB, error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
