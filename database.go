package snorlax

import (
	"database/sql"

	// sqlite3 is the database driver for sqlite.
	_ "github.com/mattn/go-sqlite3"
)

// InitDB intializes the connection to the SQLite database.
func (s *Snorlax) InitDB() {
	db, err := sql.Open("sqlite3", s.config.DBPath)
	if err != nil {
		s.Log.WithError(err).Error("Error opening connection to database.")
		return
	}

	s.DB = db
}
