package models

import (
	"database/sql"
	"errors"
)

// User holds the user representation.
// A user is server-specific. So, it isn't unique.
type User struct {
	Bans      int    // Number of times user has been banned.
	Kicks     int    // Number of times user has been kicked.
	Points    int    // Warning points.
	Portfolio string // Portfolio URL.
	ServerID  string
	UserID    string
}

// InitUsers will initialize all of the tables related to users.
func InitUsers(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `Users` (" +
		"`UserID` TEXT NOT NULL PRIMARY KEY," +
		"`ServerID` TEXT NOT NULL," +
		"`Kicks` INTEGER," +
		"`Bans` INTEGER," +
		"`Points` INTEGER," +
		"`Portfolio` TEXT" +
		")")
	if err != nil {
		return err
	}

	return nil
}

// Errors which the functions will return.
var (
	ErrUserNotExist = errors.New("user doesn't exist")
)

// GetUser will get a user with a given userID and serverID.
func GetUser(db *sql.DB, userID, serverID string) (*User, error) {
	row := db.QueryRow("SELECT * FROM Users WHERE userID=? AND serverID=?", userID, serverID)
	user := &User{}
	err := row.Scan(&user.UserID, &user.ServerID, &user.Kicks, &user.Bans, &user.Points, &user.Portfolio)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, ErrUserNotExist
	}

	return user, nil
}

// This holds the values for the different types of responses from the
// thresholds functions.
var (
	NoAction  = 0
	Kick      = 1
	BanPoints = 2
	BanKicks  = 3
)

// CheckThresholds will check if the user has hit any of the thresholds.
func (user *User) CheckThresholds(db *sql.DB) (int, error) {
	ok, err := user.CheckBanThreshold(db)
	if err != nil {
		return 0, err
	}

	if ok > 1 {
		return ok, nil
	}

	ok, err = user.CheckKickThreshold(db)
	if err != nil {
		return 0, err
	}

	return ok, nil
}

// CheckKickThreshold will check if the user has hit the kick threshold.
func (user *User) CheckKickThreshold(db *sql.DB) (int, error) {
	warnConfig, err := GetWarnConfig(db, user.ServerID)
	if err != nil {
		return 0, err
	}

	if warnConfig.KickThreshold != 0 && user.Points-(user.Kicks*warnConfig.KickThreshold) >= warnConfig.KickThreshold {
		return 1, nil
	}

	return 0, nil
}

// CheckBanThreshold will check if the user has hit the ban point or kick threshold.
func (user *User) CheckBanThreshold(db *sql.DB) (int, error) {
	warnConfig, err := GetWarnConfig(db, user.ServerID)
	if err != nil {
		return 0, err
	}

	if warnConfig.BanThreshold != 0 && user.Points-(user.Bans*warnConfig.BanThreshold) >= warnConfig.BanThreshold {
		return 2, nil
	}

	if warnConfig.BanKickThreshold != 0 && user.Kicks-(user.Bans*warnConfig.BanKickThreshold) >= warnConfig.BanKickThreshold {
		return 3, nil
	}

	return 0, nil
}

// Insert will insert the user into the database.
// Insert will update an existing user if it already exists.
func (user *User) Insert(db *sql.DB) error {
	emptyStr := ""
	err := db.QueryRow("SELECT UserID FROM Users WHERE UserID=? AND ServerID=?", user.UserID, user.ServerID).Scan(&emptyStr)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO Users (UserID, ServerID, Kicks, Bans, Points, Portfolio) values(?,?,?,?,?,?)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(user.UserID, user.ServerID, user.Kicks, user.Bans, user.Points, user.Portfolio)
		if err != nil {
			return err
		}

		return nil
	}

	stmt, err := db.Prepare("UPDATE Users SET (Kicks, Bans, Points, Portfolio) = (?,?,?,?) WHERE UserID=? AND ServerID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Kicks, user.Bans, user.Points, user.Portfolio, user.UserID, user.ServerID)
	if err != nil {
		return err
	}

	return nil
}
