package models

import (
	"database/sql"
	"errors"
)

// BirthdayInit will set up the tables required for the birthday module.
func BirthdayInit(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `Birthdays` (" +
		"`UserID` TEXT NOT NULL PRIMARY KEY," +
		"`Birthday` TEXT NOT NULL," +
		"`ServerID` TEXT NOT NULL" +
		")")

	if err != nil {
		return err
	}

	return nil
}

// ErrNoBirthdayFound is returned when there is no birthday found.
var ErrNoBirthdayFound = errors.New("no birthday was found")

// Birthday is the model.
type Birthday struct {
	UserID   string
	Birthday string // Specified as MM/DD
	ServerID string
}

// GetBirthday will check the database for a birthday.
func GetBirthday(db *sql.DB, userID, serverID string) (*Birthday, error) {
	bday := &Birthday{}
	err := db.QueryRow("SELECT * FROM Birthdays WHERE userID=? AND serverID=?", userID, serverID).Scan(&bday.UserID, &bday.Birthday, &bday.ServerID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, ErrNoBirthdayFound
	}

	return bday, nil
}

// GetBirthdaysWithDate will check if any databases have birthday at the current date.
// Birthday is of format MM/DD.
func GetBirthdaysWithDate(db *sql.DB, birthday string) ([]*Birthday, error) {
	rows, err := db.Query("SELECT * FROM Birthdays WHERE birthday=?", birthday)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, ErrNoBirthdayFound
	}

	birthdays := []*Birthday{}

	defer rows.Close()

	for rows.Next() {
		bday := &Birthday{}
		err := rows.Scan(&bday.UserID, &bday.Birthday, &bday.ServerID)
		if err != nil {
			return nil, err
		}

		birthdays = append(birthdays, bday)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return birthdays, nil
}

// Insert will insert the birthday into the database.
// Insert will override an old birthday.
func (b *Birthday) Insert(db *sql.DB) error {
	emptyStr := ""
	err := db.QueryRow("SELECT UserID FROM Birthdays WHERE userID=? AND serverID=?", b.UserID, b.ServerID).Scan(&emptyStr)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO Birthdays (UserID, Birthday, ServerID) values(?,?,?)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(b.UserID, b.Birthday, b.ServerID)
		if err != nil {
			return err
		}

		return nil
	}

	stmt, err := db.Prepare("UPDATE Birthdays SET (Birthday) = (?) WHERE userID=? AND serverID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(b.Birthday, b.UserID, b.ServerID)
	if err != nil {
		return err
	}

	return nil
}
