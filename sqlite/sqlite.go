package sqlite

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/luizbranco/waukeen"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func New(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS accounts(
			id INTEGER PRIMARY KEY,
			number TEXT NOT NULL CHECK(number <> ''),
			name TEXT,
			type INTEGER NOT NULL,
			currency TEXT,
			balance INTEGER
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS transactions(
			id INTEGER PRIMARY KEY,
			account_id INTEGER,
			fitid TEXT NOT NULL,
			type INTEGER NOT NULL,
			title TEXT NOT NULL,
			alias TEXT,
			description TEXT,
			amount INTEGER,
			date DATETIME,
			FOREIGN KEY(account_id) REFERENCES accounts(id)
		);
		`,
		`
		CREATE UNIQUE INDEX IF NOT EXISTS account_fitid ON
		transactions(account_id, fitid)
		`,
		`
		CREATE TABLE IF NOT EXISTS tags(
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS transaction_tags(
			transaction_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			FOREIGN KEY(transaction_id) REFERENCES transactions(id)
			FOREIGN KEY(tag_id) REFERENCES tags(id)
		);
		`,
		`
		CREATE UNIQUE INDEX IF NOT EXISTS transaction_tag ON
		transaction_tags(transaction_id, tag_id)
		`,
		`
		CREATE TABLE IF NOT EXISTS rules(
			id INTEGER PRIMARY KEY,
			account_id INTEGER,
			type INTEGER NOT NULL,
			match TEXT NOT NULL,
			result TEXT NOT NULL,
			FOREIGN KEY(account_id) REFERENCES accounts(id)
		);
		`,
	}

	for _, q := range queries {
		_, err = db.Exec(q)

		if err != nil {
			return nil, err
		}
	}

	return &DB{db}, nil
}

func (db *DB) CreateAccount(a *waukeen.Account) error {
	q := `INSERT into accounts (number, name, type, currency, balance) values
(?, ?, ?, ?, ?);`

	res, err := db.Exec(q, a.Number, a.Name, a.Type, a.Currency, a.Balance)

	if err != nil {
		return fmt.Errorf("error creating account: %s", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return fmt.Errorf("error retrieving last account id: %s", err)
	}

	a.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) FindAccounts() ([]waukeen.Account, error) {
	var accounts []waukeen.Account

	rows, err := db.Query("SELECT id, number, name, type, currency, balance FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		a := waukeen.Account{}
		err = rows.Scan(&a.ID, &a.Number, &a.Name, &a.Type, &a.Currency, &a.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	err = rows.Err()
	return accounts, err
}

func (db *DB) FindAccount(number string) (*waukeen.Account, error) {
	q := "SELECT id, number, name, type, currency, balance FROM accounts where number = ?"

	a := &waukeen.Account{}

	err := db.QueryRow(q, number).Scan(&a.ID, &a.Number, &a.Name, &a.Type,
		&a.Currency, &a.Balance)

	if err != nil {
		return nil, fmt.Errorf("error finding account: %s", err)
	}

	return a, nil
}

func (db *DB) UpdateAccount(a *waukeen.Account) error {
	_, err := db.Exec(`
	UPDATE accounts SET number=?, name=?, type=?, currency=?, balance=?
	where id = ?`, a.Number, a.Name, a.Type, a.Currency, a.Balance, a.ID)
	return err
}

func (db *DB) CreateTransaction(t *waukeen.Transaction) error {
	q := `INSERT OR IGNORE into transactions
	(account_id, fitid, type, title, alias, description, amount, date)
	values (?, ?, ?, ?, ?, ?, ?, ?);`

	res, err := db.Exec(q, t.AccountID, t.FITID, t.Type, t.Title, t.Alias,
		t.Description, t.Amount, t.Date)

	if err != nil {
		return fmt.Errorf("error creating transaction: %s", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return fmt.Errorf("error retrieving last transaction id: %s", err)
	}

	t.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) FindTransactions(opts waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error) {
	var transactions []waukeen.Transaction
	var clauses []string

	if len(opts.Accounts) > 0 {
		clause := "account_id IN (" + strings.Join(opts.Accounts, " ,") + ")"
		clauses = append(clauses, clause)
	}

	if len(opts.Types) > 0 {
		var types []string

		for _, t := range opts.Types {
			types = append(types, strconv.Itoa(int(t)))
		}

		clause := "type IN (" + strings.Join(types, " ,") + ")"
		clauses = append(clauses, clause)
	}

	if len(opts.Tags) > 0 {
		//FIXME
	}

	if !opts.Start.IsZero() {
		clauses = append(clauses, "date >= "+opts.Start.Format("'2006-01-02'"))
	}

	if !opts.End.IsZero() {
		clauses = append(clauses, "date <= "+opts.End.Format("'2006-01-02'"))
	}

	q := `SELECT id, account_id, type, title, alias, description, amount, date
	FROM transactions WHERE `

	q += strings.Join(clauses, " AND ")

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := waukeen.Transaction{}
		err = rows.Scan(&t.ID, &t.AccountID, &t.Type, &t.Title, &t.Alias,
			&t.Description, &t.Amount, &t.Date)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	err = rows.Err()
	return transactions, err
}

func (db *DB) CreateRule(r *waukeen.Rule) error {
	q := `INSERT into rules (account_id, type, match, result) values
		(?, ?, ?, ?);`

	res, err := db.Exec(q, r.AccountID, r.Type, r.Match, r.Result)

	if err != nil {
		return fmt.Errorf("error creating rule: %s", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return fmt.Errorf("error retrieving last rule id: %s", err)
	}

	r.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) FindRules(acc string) ([]waukeen.Rule, error) {
	var rules []waukeen.Rule

	stmt := `SELECT id, account_id, type, match, result from rules
	where account_id = ? OR account_id = ""`

	rows, err := db.Query(stmt, acc)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := waukeen.Rule{}
		err = rows.Scan(&r.ID, &r.AccountID, &r.Type, &r.Match, &r.Result)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	err = rows.Err()
	return rules, err
}

func (db *DB) CreateStatement(stmt waukeen.Statement,
	transformer waukeen.TransactionTransformer) error {
	number := stmt.Account.Number

	acc, err := db.FindAccount(number)

	if err == nil {
		acc.Balance = stmt.Account.Balance
		err = db.UpdateAccount(acc)
	} else {
		acc = &stmt.Account
		err = db.CreateAccount(acc)
	}

	if err != nil {
		return err
	}

	rules, err := db.FindRules(acc.ID)

	if err != nil {
		return err
	}

	for _, tn := range stmt.Transactions {
		t := &tn
		t.AccountID = acc.ID
		for _, r := range rules {
			transformer.Transform(t, r)
		}
		err := db.CreateTransaction(t)
		if err != nil {
			return err
		}
	}

	return nil
}
