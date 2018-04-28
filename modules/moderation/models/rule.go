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

// GetRule will get a rule with a specified ID.
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

// DelRule will delete a rule with a specified ID.
func DelRule(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM Rules WHERE ID=?", id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		return ErrRuleNotExist
	}

	return nil
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

		res, err := stmt.Exec(rule.Points, rule.Description, rule.ServerID)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}

		rule.ID = int(id)

		return nil
	}

	stmt, err := db.Prepare("UPDATE Rules SET (Points, Description) = (?,?) WHERE ID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(rule.Points, rule.Description, rule.ID)
	if err != nil {
		return err
	}

	return nil
}

// serverRulesCache is the cached versions of the different server rules.
// It is a map of serverID to *ServerRules.
var serverRulesCache = map[string]*ServerRules{}

// GetServerRules will get a list of server rules with a serverID.
func GetServerRules(db *sql.DB, serverID string) (*ServerRules, error) {
	serverRules, ok := serverRulesCache[serverID]
	if ok {
		return serverRules, nil
	}

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
		if ruleListStr[i] == "" {
			break
		}

		ruleID, err := strconv.Atoi(ruleListStr[i])
		if err != nil {
			return nil, err
		}

		ruleListInt = append(ruleListInt, ruleID)
	}

	serverRules = &ServerRules{
		ServerID: serverID,
		RuleIDs:  ruleListInt,
	}

	serverRulesCache[serverID] = serverRules
	return serverRules, nil
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

// DelRule will remove a rule from the list of rules.
func (serverRules *ServerRules) DelRule(db *sql.DB, ruleIdx int) error {
	err := serverRules.UpdateServerRules(db)
	if err != nil {
		return err
	}

	err = DelRule(db, serverRules.RuleIDs[ruleIdx])
	if err != nil {
		return err
	}
	serverRules.RuleIDs = append(serverRules.RuleIDs[:ruleIdx], serverRules.RuleIDs[ruleIdx+1:]...)

	err = serverRules.Insert(db)
	if err != nil {
		return err
	}

	return nil
}

// Insert will insert the server rules into the database.
// Insert will override old serv er rules.
func (serverRules *ServerRules) Insert(db *sql.DB) error {
	emptyStr := ""
	err := db.QueryRow("SELECT ServerID FROM ServerRules WHERE ServerID=?", serverRules.ServerID).Scan(&emptyStr)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	serverRulesStr := ""
	for i := 0; i < len(serverRules.RuleIDs); i++ {
		if i == len(serverRules.RuleIDs)-1 {
			serverRulesStr += strconv.Itoa(serverRules.RuleIDs[i])
		} else {
			serverRulesStr += strconv.Itoa(serverRules.RuleIDs[i]) + ","
		}
	}

	if err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO ServerRules (ServerID, RuleIDs) values(?,?)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(serverRules.ServerID, serverRulesStr)
		if err != nil {
			return err
		}

		serverRulesCache[serverRules.ServerID] = serverRules
		return nil
	}

	stmt, err := db.Prepare("UPDATE ServerRules SET (RuleIDs) = (?) WHERE ServerID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(serverRulesStr, serverRules.ServerID)
	if err != nil {
		return err
	}

	serverRulesCache[serverRules.ServerID] = serverRules
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
		if ruleListStr[i] == "" {
			break
		}

		ruleID, err := strconv.Atoi(ruleListStr[i])
		if err != nil {
			return err
		}

		ruleListInt = append(ruleListInt, ruleID)
	}

	serverRules.RuleIDs = ruleListInt
	serverRulesCache[serverRules.ServerID] = serverRules
	return nil
}
