package models

import (
	"database/sql"
)

// Init initiliazes all database tables.
func Init(db *sql.DB) error {
	err := InitCurrency(db)
	if err != nil {
		return err
	}

	return nil
}
