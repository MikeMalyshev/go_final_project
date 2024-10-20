package db

import (
	"database/sql"
	"fmt"
	"os"

	"go_final_project/internal/config"

	stdlib "github.com/multiprocessio/go-sqlite3-stdlib"
)

var sqlite3_ext_registred bool

type DBStorage struct {
	db  *sql.DB
	cfg *config.Handler
}

func New(cfg *config.Handler) *DBStorage {
	if !sqlite3_ext_registred {
		stdlib.Register("sqlite3_ext")
		sqlite3_ext_registred = true
	}
	return &DBStorage{cfg: cfg}
}

func (storage *DBStorage) Create() error {
	dbPath := storage.cfg.DBPath()
	if dbPath == "" {
		return fmt.Errorf("DBPath have not been setted ")
	}

	if storage.Exists() {
		return fmt.Errorf("%s already exists ", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening sqlite: %v ", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE scheduler (
		id 		INTEGER 		PRIMARY KEY AUTOINCREMENT,
		date 	CHAR(8) 		NOT NULL,
		title 	VARCHAR(128) 	NOT NULL,
		comment TEXT 			NOT NULL 	DEFAULT "",
		repeat 	VARCHAR(128) 	NOT NULL 	DEFAULT ""
		)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX date_index ON scheduler (date)`)
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) Exists() bool {
	if _, err := os.Stat(storage.cfg.DBPath()); err == nil {
		return true
	}
	return false
}

func (storage *DBStorage) Open() error {
	var err error
	storage.db, err = sql.Open("sqlite3_ext", storage.cfg.DBPath())
	if err != nil {
		return fmt.Errorf("DBStorage.Open: %v", err)
	}
	return nil
}

func (storage *DBStorage) Close() error {
	return storage.db.Close()
}
