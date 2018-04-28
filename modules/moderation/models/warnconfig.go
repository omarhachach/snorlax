package models

import (
	"database/sql"
	"errors"
)

// WarnConfig is the server configuration for the warn/kick/ban logging.
type WarnConfig struct {
	ServerID     string
	LogChannelID string
	LogWarn      bool
	LogKick      bool
	LogBan       bool
}

// InitWarnConfig will intialize the tables for WarnConfigs.
func InitWarnConfig(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `WarnConfigs` (" +
		"`ServerID` TEXT NOT NULL PRIMARY KEY UNIQUE," +
		"`LogChannelID` TEXT," +
		"`LogWarn` BOOLEAN," +
		"`LogKick` BOOLEAN," +
		"`LogBan` BOOLEAN" +
		")")
	if err != nil {
		return err
	}

	return nil
}

// This holds the errors which the functions will return.
var (
	ErrWarnConfigNotExist = errors.New("specified warn config doesn't exist")
)

// warnConfigCache holds the cached versions of the WarnConfigs.
// It is mapped as serverID to *WarnConfig.
var warnConfigCache = map[string]*WarnConfig{}

// GetWarnConfig will retrieve the WarnConfig for a specific server.
func GetWarnConfig(db *sql.DB, serverID string) (*WarnConfig, error) {
	warnConfig, ok := warnConfigCache[serverID]
	if ok {
		return warnConfig, nil
	}

	row := db.QueryRow("SELECT * FROM WarnConfigs WHERE ServerID=?", serverID)

	warnConfig = &WarnConfig{}

	err := row.Scan(&warnConfig.ServerID, &warnConfig.LogChannelID, &warnConfig.LogWarn, &warnConfig.LogKick, &warnConfig.LogBan)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, ErrWarnConfigNotExist
	}

	warnConfigCache[serverID] = warnConfig
	return warnConfig, nil
}

// Insert will insert the server rules into the database.
// Insert will override old serv er rules.
func (warnConfig *WarnConfig) Insert(db *sql.DB) error {
	emptyStr := ""
	err := db.QueryRow("SELECT ServerID FROM WarnConfigs WHERE ServerID=?", warnConfig.ServerID).Scan(&emptyStr)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO WarnConfigs (ServerID, LogChannelID, LogWarn, LogKick, LogBan) values(?,?,?,?,?)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(warnConfig.ServerID, warnConfig.LogChannelID, warnConfig.LogWarn, warnConfig.LogKick, warnConfig.LogBan)
		if err != nil {
			return err
		}

		warnConfigCache[warnConfig.ServerID] = warnConfig
		return nil
	}

	stmt, err := db.Prepare("UPDATE ServerRules SET (LogChannelID, LogWarn, LogKick, LogBan) = (?,?,?,?) WHERE ServerID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(warnConfig.LogChannelID, warnConfig.LogWarn, warnConfig.LogKick, warnConfig.LogBan, warnConfig.ServerID)
	if err != nil {
		return err
	}

	warnConfigCache[warnConfig.ServerID] = warnConfig
	return nil
}
