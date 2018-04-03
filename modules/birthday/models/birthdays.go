package models

import (
	"errors"
	"strings"

	"database/sql"
)

// CurrentBirthday holds all of the ServerIDs where a given user has gotten a
// birthday role.
type CurrentBirthday struct {
	UserID          string
	ServerIDs       []string
	Birthday        string
	BirthdayRoleIDs []string
}

// CurrentBirthdaysInit will initialize the CurrentBirthdays table.
func CurrentBirthdaysInit(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `CurrentBirthdays` (" +
		"`UserID` TEXT NOT NULL PRIMARY KEY UNIQUE," +
		"`ServerIDS` TEXT NOT NULL," +
		"`Birthday` TEXT NOT NULL," +
		"`BirthdayRoleID` TEXT NOT NULL" +
		")")

	if err != nil {
		return err
	}

	return nil
}

// ErrNoCurrentBirthdayFound is the error value for not finding any current
// birthdays.
var ErrNoCurrentBirthdayFound = errors.New("no current birthdays were found")

// GetCurrentBirthdays will retrieve all of the current birthdays.
func GetCurrentBirthdays(db *sql.DB) ([]*CurrentBirthday, error) {
	rows, err := db.Query("SELECT * FROM CurrentBirthdays")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, ErrNoCurrentBirthdayFound
	}

	currBdays := []*CurrentBirthday{}

	defer rows.Close()

	for rows.Next() {
		currBday := &CurrentBirthday{}
		serverIDs := ""
		birthdaysRoleIDs := ""
		err := rows.Scan(&currBday.UserID, &serverIDs, &currBday.Birthday, &birthdaysRoleIDs)
		if err != nil {
			return nil, err
		}

		currBday.ServerIDs = strings.Split(serverIDs, ",")
		currBday.BirthdayRoleIDs = strings.Split(birthdaysRoleIDs, ",")

		currBdays = append(currBdays, currBday)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return currBdays, nil
}

// DeleteCurrentBirthday will remove a given user from the CurrentBirthdays
// table.
func DeleteCurrentBirthday(db *sql.DB, userID string) error {
	_, err := db.Exec("DELETE FROM CurrentBirthdays WHERE userID=? LIMIT 1", userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return ErrNoBirthdayConfigFound
	}

	return nil
}
