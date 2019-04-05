package wikipedia

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	log.Printf("Query: %s\nParams: %v", query, args)
	return db.DB.QueryRow(query, args...)
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	log.Printf("Query: %s\nParams: %v", query, args)
	return db.DB.Exec(query, args...)
}

func OpenDB(db_url string) DB {
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	return DB{db}
}

func CloseDB(db DB) {
	db.Close()
}
