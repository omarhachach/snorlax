package models

import (
	"database/sql"
	"errors"
)

// BirthdayConfig holds the birthday config for a server.
type BirthdayConfig struct {
	ServerID       string
	AssignRole     bool
	BirthdayRoleID string
}

// BirthdayConfigInit initializes the birthday config table.
func BirthdayConfigInit(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `BirthdayConfigs` (" +
		"`ServerID` TEXT NOT NULL PRIMARY KEY UNIQUE," +
		"`AssignRole` BOOLEAN NOT NULL," +
		"`BirthdayRoleID` TEXT" +
		")")

	if err != nil {
		return err
	}

	return nil
}

// ErrNoBirthdayConfigFound is the error for not finding a birthday config.
var ErrNoBirthdayConfigFound = errors.New("no birthday config was found")

// GetBirthdayConfig returns a specific
func GetBirthdayConfig(db *sql.DB, serverID string) (*BirthdayConfig, error) {
	bdayConfig := &BirthdayConfig{
		ServerID: serverID,
	}

	err := db.QueryRow("SELECT AssignRole,BirthdayRoleID FROM BirthdayConfigs WHERE serverID=?", bdayConfig.ServerID).Scan(&bdayConfig.AssignRole, &bdayConfig.BirthdayRoleID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, ErrNoBirthdayConfigFound
	}

	return bdayConfig, nil
}

// Insert will insert the birthday config into the database.
// Insert will override the old birthday config.
func (bdayConfig *BirthdayConfig) Insert(db *sql.DB) error {
	emptyStr := ""
	err := db.QueryRow("SELECT ServerID FROM BirthdayConfigs WHERE serverID=?", bdayConfig.ServerID).Scan(&emptyStr)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO BirthdayConfigs (ServerID, AssignRole, BirthdayRoleID) values(?,?,?)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(bdayConfig.ServerID, bdayConfig.AssignRole, bdayConfig.BirthdayRoleID)
		if err != nil {
			return err
		}

		return nil
	}

	stmt, err := db.Prepare("UPDATE BirthdayConfigs SET (AssignRole, BirthdayRoleID) = (?,?) WHERE serverID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(bdayConfig.AssignRole, bdayConfig.BirthdayRoleID, bdayConfig.ServerID)
	if err != nil {
		return err
	}

	return nil
}
