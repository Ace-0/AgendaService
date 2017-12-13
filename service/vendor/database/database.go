package database

import (
	"database/sql"
	"log"
	"os"
)

const createUserTableSQL = `CREATE TABLE IF NOT EXISTS User (
	username TEXT PRIMARY KEY,
	password TEXT NOT NULL,
	email TEXT NOT NULL,
	phone TEXT NOT NULL
);`

const createLoginTableSQL = `CREATE TABLE IF NOT EXISTS Login (
	token TEXT NOT NULL PRIMARY KEY,
	username TEXT NOT NULL,
	FOREIGN KEY(username) REFERENCES User(username) ON DELETE CASCADE
);`

type sqlWork func(db *sql.DB)

type agendaDB struct {
	db    *sql.DB
	file  string
	queue chan sqlWork
}

var adb *agendaDB

// WithDB ..
func WithDB(file string, f func()) {
	InitializeDB(file)
	log.Print("database set up")
	f()
	ClearDB()
	log.Print("database tear down")
}

// WithTestDB ..
func WithTestDB(f func()) {
	WithDB("Test.db", f)
}

// CloseDB do not support concurrent access
func CloseDB() error {
	if adb == nil {
		return nil
	}
	var err error
	if adb.db != nil {
		err = adb.db.Close()
	}
	if adb.db != nil {
		close(adb.queue)
	}
	adb = nil
	return err
}

func prepareDB(dbfile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbfile)
	// open database fail
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	_, err = db.Exec(createUserTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = db.Exec(createLoginTableSQL)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	log.Print("database initialized")
	return db, nil
}

func startExecQueue(db *sql.DB) chan sqlWork {
	q := make(chan sqlWork)
	go func() {
		for f := range q {
			f(db)
		}
	}()
	return q
}

// InitializeDB with specific file
func InitializeDB(dbfile string) error {
	CloseDB()
	db, err := prepareDB(dbfile)
	if err != nil {
		return err
	}
	adb = &agendaDB{
		db:    db,
		file:  dbfile,
		queue: startExecQueue(db),
	}
	return nil
}

// ClearDB and the file
func ClearDB() error {
	if adb == nil {
		return nil
	}
	os.Remove(adb.file)
	return CloseDB()
}

func pQuery(rf func(*sql.Rows), ef func(error), s string, args ...interface{}) chan struct{} {
	done := make(chan struct{})
	adb.queue <- func(db *sql.DB) {
		defer func() { done <- struct{}{} }()
		row, err := db.Query(s, args...)
		if err != nil {
			ef(err)
			return
		}
		for row.Next() {
			if row.Err() != nil {
				ef(row.Err())
				return
			}
			rf(row)
		}
		row.Close()
	}
	return done
}

func pExec(s string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var e error
	done := make(chan struct{})
	adb.queue <- func(db *sql.DB) {
		res, err := db.Exec(s, args...)
		result = res
		e = err
		done <- struct{}{}
	}
	<-done
	return result, e
}
