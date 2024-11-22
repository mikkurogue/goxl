package database

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Store     *sql.DB
	Connected bool
}

func (db *Database) Connect() error {

	store, err := sql.Open("sqlite3", "./store.db")
	if err != nil {
		return err
	}

	db.Store = store
	db.Connected = true

	return nil
}

func (db *Database) Disconnect() error {
	if db.Connected == false {
		return errors.New("Connection never established")
	}

	db.Connected = false
	defer db.Store.Close()
	return nil
}

func (db *Database) CreateTable() error {

	if db.Store == nil {
		return errors.New("No database connection")
	}

	return nil
}
