package snorlax

import (
	"database/sql"

	// sqlite3 is the database driver for sqlite.
	_ "github.com/mattn/go-sqlite3"
)

// ConnDB intializes the connection to the SQLite database.
func (s *Snorlax) ConnDB() {
	db, err := sql.Open("sqlite3", s.config.DBPath)
	if err != nil {
		s.Log.WithError(err).Error("Error opening connection to database.")
		return
	}

	err = db.Ping()
	if err != nil {
		s.Log.WithError(err).Error("Error pinging database connection.")
		return
	}

	s.DB = db
}

// InitDB will create the initial tables and rows of the database.
func (s *Snorlax) InitDB() {

}
