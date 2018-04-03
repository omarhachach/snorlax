package models

import (
	"database/sql"
	"time"
)

// InitCurrency initiliazes the tables needed for the currency system.
func InitCurrency(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `Currency` (" +
		"`UserID` TEXT NOT NULL PRIMARY KEY UNIQUE," +
		"`Amount` INTEGER NOT NULL" +
		")")

	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `CurrencyTransactions` (" +
		"`ID` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT," +
		"`Amount` INTEGER NOT NULL," +
		"`Reason` TEXT NOT NULL," +
		"`UserID` TEXT NOT NULL," +
		"`DateAdded` INTEGER NOT NULL" +
		")")

	if err != nil {
		return err
	}

	return nil
}

// Currency is the type for holding currency info about a user.
type Currency struct {
	UserID    string
	Amount    int
	DateAdded int
}

// CreateCurrency creates and inserts a Currency.
func CreateCurrency(db *sql.DB, amount int, userID string) (*Currency, error) {
	stmt, err := db.Prepare("INSERT INTO Currency (Amount, UserID, DateAdded) values(?,?,?)")
	if err != nil {
		return nil, err
	}

	dateAdded := int(time.Now().Unix())
	_, err = stmt.Exec(amount, userID, dateAdded)
	if err != nil {
		return nil, err
	}

	return &Currency{
		UserID:    userID,
		Amount:    amount,
		DateAdded: dateAdded,
	}, nil
}

// GetCurrency returns a specific currency
func GetCurrency(db *sql.DB, userID string) (*Currency, error) {
	currency := &Currency{}
	err := db.QueryRow("SELECT * FROM Currency WHERE userID =?", userID).Scan(&currency.UserID, &currency.Amount, &currency.DateAdded)
	if err != nil {
		return nil, err
	}

	return currency, nil
}

// CurrencyTransaction is the type for keeping track of a transaction to a
// currency type.
type CurrencyTransaction struct {
	ID        int
	Amount    int
	UserID    string
	Reason    string
	DateAdded int
}

// CreateCurrencyTransaction will create and insert a new CurrencyTransaction.
func CreateCurrencyTransaction(db *sql.DB, amount int, userID, reason string) (*CurrencyTransaction, error) {
	stmt, err := db.Prepare("INSERT INTO CurrencyTransaction (Amount, UserID, Reason, DateAdded) values (?,?,?,?)")
	if err != nil {
		return nil, err
	}

	dateAdded := int(time.Now().Unix())
	_, err = stmt.Exec(amount, userID, reason, dateAdded)
	if err != nil {
		return nil, err
	}

	id := 0
	err = db.QueryRow("SELECT ID FROM CurrencyTransaction WHERE UserID =? AND WHERE DateAdded =?", userID, dateAdded).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &CurrencyTransaction{
		ID:        id,
		Amount:    amount,
		UserID:    userID,
		Reason:    reason,
		DateAdded: dateAdded,
	}, nil
}
