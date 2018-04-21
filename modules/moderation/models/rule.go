package models

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

// ServerRules holds a list of rule IDs (IDs in the database.).
type ServerRules struct {
	ServerID string
	RuleIDs  []int // Represented in database as: ruleID,ruleID,ruleID...
}

// Rule holds a rule and its points.
type Rule struct {
	ID          int
	Points      int
	ServerID    string
	Description string
}

// InitRule will initialize the rules database.
func InitRule(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `ServerRules` (" +
		"`ServerID` TEXT NOT NULL PRIMARY KEY," +
		"`RuleIDs` TEXT" +
		")")

	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `Rules` (" +
		"`ID` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`Points` INTEGER NOT NULL," +
		"`ServerID` TEXT NOT NULL," +
		"`Description` TEXT NOT NULL" +
		")")

	if err != nil {
		return err
	}

	return nil
}

// All of the errors which the queries will return.
var (
	ErrRuleNotExist         = errors.New("specified rule does not exist")
	ErrServerDoesntHaveRule = errors.New("specified server rule ID isn't in the rule list")
	ErrServerRulesDontExist = errors.New("specified server rule list does not exist")
)

// GetRule will get a rule with a specifeid ID.
func GetRule(db *sql.DB, id int) (*Rule, error) {
	row := db.QueryRow("SELECT * FROM Rules WHERE id=?", id)
	rule := &Rule{}
	err := row.Scan(&rule.ID, &rule.Points, &rule.ServerID, &rule.Description)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, ErrRuleNotExist
	}

	return rule, nil
}

// Insert will insert the rule into the database.
// Insert will override an old rule.
func (rule *Rule) Insert(db *sql.DB) error {
	emptyStr := ""
	err := db.QueryRow("SELECT ID FROM Rules WHERE ID=?", rule.ID).Scan(&emptyStr)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO Rules (Points, Description, ServerID) values(?,?,?)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(rule.Points, rule.Description, rule.ServerID)
		if err != nil {
			return err
		}

		return nil
	}

	stmt, err := db.Prepare("UPDATE Rules SET (Points, Description) = (?,?) WHERE ID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(rule.Points, rule.Description)
	if err != nil {
		return err
	}

	return nil
}

// GetServerRules will get a list of server rules with a serverID.
func GetServerRules(db *sql.DB, serverID string) (*ServerRules, error) {
	rulesStr := ""
	err := db.QueryRow("SELECT RuleIDs FROM ServerRules WHERE serverID=?", serverID).Scan(&rulesStr)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, ErrServerRulesDontExist
	}

	ruleListStr := strings.Split(rulesStr, ",")
	ruleListInt := make([]int, 0, len(ruleListStr))
	for i := 0; i < len(ruleListStr); i++ {
		ruleID, err := strconv.Atoi(ruleListStr[i])
		if err != nil {
			return nil, err
		}

		ruleListInt = append(ruleListInt, ruleID)
	}

	return &ServerRules{
		ServerID: serverID,
		RuleIDs:  ruleListInt,
	}, nil
}

// AddRule will add a rule to the list of rules.
func (serverRules *ServerRules) AddRule(db *sql.DB, rule *Rule) error {
	err := serverRules.UpdateServerRules(db)
	if err != nil {
		return err
	}

	serverRules.RuleIDs = append(serverRules.RuleIDs, rule.ID)
	err = serverRules.Insert(db)
	if err != nil {
		return err
	}

	return nil
}

// DelRule will remove a rule from the list of rules-
func (serverRules *ServerRules) DelRule(db *sql.DB, ruleID int) error {
	err := serverRules.UpdateServerRules(db)
	if err != nil {
		return err
	}

	for i := 0; i < len(serverRules.RuleIDs); i++ {
		if serverRules.RuleIDs[i] == ruleID {
			serverRules.RuleIDs = append(serverRules.RuleIDs[:i], serverRules.RuleIDs[i+1:]...)
			return nil
		}
	}

	return ErrServerDoesntHaveRule
}

// Insert will insert the server rules into the database.
// Insert will override old serv er rules.
func (serverRules *ServerRules) Insert(db *sql.DB) error {
	emptyStr := ""
	err := db.QueryRow("SELECT ServerID FROM ServerRules WHERE ServerID=?", serverRules.ServerID).Scan(&emptyStr)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO ServerRules (ServerID, RuleIDs) values(?,?)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(serverRules.ServerID, serverRules.RuleIDs)
		if err != nil {
			return err
		}

		return nil
	}

	stmt, err := db.Prepare("UPDATE ServerRules SET (RuleIDs) = (?) WHERE ServerID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(serverRules.RuleIDs, serverRules.ServerID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateServerRules will update the serverRules pointer information.
func (serverRules *ServerRules) UpdateServerRules(db *sql.DB) error {
	rulesStr := ""
	err := db.QueryRow("SELECT RuleIDs FROM ServerRules WHERE ServerID=?", serverRules.ServerID).Scan(&rulesStr)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return ErrServerRulesDontExist
	}

	ruleListStr := strings.Split(rulesStr, ",")
	ruleListInt := make([]int, 0, len(ruleListStr))
	for i := 0; i < len(ruleListStr); i++ {
		ruleID, err := strconv.Atoi(ruleListStr[i])
		if err != nil {
			return err
		}

		ruleListInt = append(ruleListInt, ruleID)
	}

	serverRules.RuleIDs = ruleListInt
	return nil
}
