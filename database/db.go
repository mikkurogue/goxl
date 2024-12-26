package database

import (
	"database/sql"
	"errors"
	"fmt"
	"goxl/util"
	"os"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Store     *sql.DB
	Connected bool
}

type ProcessRow struct {
	Id       string `json:"id"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

func checkIfStoreExist() bool {
	if _, err := os.Stat("./store.db"); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
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

func (db *Database) CreateTable(config *util.Config) error {
	if checkIfStoreExist() {
		return nil
	}

	if db.Store == nil {
		return errors.New("No database connection")
	}

	if config == nil {
		panic("NO CONFIG PASSED TO TABLE CREATION")
	}

	fmt.Println(columnsToQuery(config.Columns))

	_, err := db.Store.Exec(columnsToQuery(config.Columns))
	if err != nil {
		panic(err)
	}

	color.Green("Succesfully created the table.")

	return nil
}

func columnsToQuery(columns []util.Column) string {
	var stmt string = "create table uploads ("

	for i := 0; i < len(columns); i++ {
		if columns[i].IsPrimaryKey {
			stmt += fmt.Sprintf(
				"%v %v primary key",
				columns[i].ColumnName,
				columns[i].Type,
			)
		} else {
			stmt += fmt.Sprintf(
				"%v %v",
				columns[i].ColumnName,
				columns[i].Type,
			)
		}

		// Add a comma unless it's the last column
		if i != len(columns)-1 {
			stmt += ", "
		}
	}

	stmt += ");"

	return stmt
}

func (db *Database) AddUpload(row ProcessRow) error {
	if db.Store == nil {
		panic("No database connection")
	}

	tx, err := db.Store.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into uploads (id, file_name, file_size) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(row.Id, row.FileName, row.FileSize)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetProcess(processId string) (ProcessRow, error) {
	if !checkIfStoreExist() {
		return ProcessRow{}, errors.New("No database connection")
	}

	rows, err := db.Store.Query("select * from uploads where id = ?", processId)
	if err != nil {
		return ProcessRow{}, err
	}
	defer rows.Close()

	var p ProcessRow

	if rows.Next() {
		err := rows.Scan(&p.Id, &p.FileName, &p.FileSize)
		if err != nil {
			return ProcessRow{}, err
		}
	} else {
		return ProcessRow{}, errors.New("no process for given id found")
	}

	return p, nil
}

func (db *Database) InsertRow(columns, values []string) error {
	if !checkIfStoreExist() {
		return errors.New("No database connection")
	}

	tx, err := db.Store.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(columnsToInsertTx(columns))
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := make([]interface{}, len(values))
	for i, v := range values {
		args[i] = v
	}

	_, execErr := stmt.Exec(args...)
	if execErr != nil {
		tx.Rollback()
		return execErr
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return commitErr
	}

	return nil
}

func columnsToInsertTx(columns []string) string {
	var stmt string = "insert into uploads ("

	for i := 0; i < len(columns); i++ {

		stmt += columns[i]

		if i != len(columns)-1 {
			stmt += ", "
		}
	}

	stmt += ") values ("

	for i := 0; i < len(columns); i++ {
		stmt += "?"

		if i != len(columns)-1 {
			stmt += ", "
		}
	}

	stmt += ");"

	return stmt
}
