package models

import (
	"database/sql"
	"errors"
)

// User holds the user representation.
// A user is server-specific. So, it isn't unique.
type User struct {
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
	err := row.Scan(&user.UserID, &user.ServerID, &user.Kicks, &user.Points, &user.Portfolio)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, ErrUserNotExist
	}

	return user, nil
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
		stmt, err := db.Prepare("INSERT INTO Users (UserID, ServerID, Kicks, Points, Portfolio) values(?,?,?,?,?)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(user.UserID, user.ServerID, user.Kicks, user.Points, user.Portfolio)
		if err != nil {
			return err
		}

		return nil
	}

	stmt, err := db.Prepare("UPDATE Users SET (Kicks, Points, Portfolio) = (?,?,?) WHERE UserID=? AND ServerID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Kicks, user.Points, user.Portfolio, user.UserID, user.ServerID)
	if err != nil {
		return err
	}

	return nil
}
