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

type Accounts DB
type Transactions DB
type Rules DB

func New(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", "./waukeen.db")
	if err != nil {
		return nil, err
	}

	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS accounts(
			id INTEGER PRIMARY KEY,
			number TEXT NOT NULL,
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
			description TEXT,
			amount INTEGER,
			date TEXT,
			tags TEXT,
			FOREIGN KEY(account_id) REFERENCES accounts(id)
		);
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

func (db *DB) Accounts() *Accounts {
	return &Accounts{db.DB}
}

func (db *DB) Transactions() *Transactions {
	return &Transactions{db.DB}
}

func (db *DB) Rules() *Rules {
	return &Rules{db.DB}
}

func (db *Accounts) Create(a *waukeen.Account) error {
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

func (db *Accounts) FindAll() ([]waukeen.Account, error) {
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

func (db *Accounts) Find(number string) (*waukeen.Account, error) {
	q := "SELECT id, number, name, type, currency, balance FROM accounts where number = ?"

	a := &waukeen.Account{}

	err := db.QueryRow(q, number).Scan(&a.ID, &a.Number, &a.Name, &a.Type,
		&a.Currency, &a.Balance)

	if err != nil {
		return nil, fmt.Errorf("error finding account: %s", err)
	}

	return a, nil
}

func (db *Accounts) Update(a *waukeen.Account) error {
	_, err := db.Exec(`
	UPDATE accounts SET number=?, name=?, type=?, currency=?, balance=?
	where id = ?`, a.Number, a.Name, a.Type, a.Currency, a.Balance, a.ID)
	return err
}

func (db *Transactions) Create(t *waukeen.Transaction) error {
	q := `INSERT into transactions
	(account_id, fitid, type, title, description, amount, tags, date)
	values (?, ?, ?, ?, ?, ?, ?, ?);`

	tags := strings.Join(t.Tags, ",")

	res, err := db.Exec(q, t.AccountID, t.FITID, t.Type, t.Title, t.Description,
		t.Amount, tags, t.Date)

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

func (db *Transactions) FindAll(acc string) ([]waukeen.Transaction, error) {
	var transactions []waukeen.Transaction

	rows, err := db.Query(`SELECT id, account_id, type, title, description,
	amount, date, tags FROM transactions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tags string
		var date string //FIXME

		t := waukeen.Transaction{}
		err = rows.Scan(&t.ID, &t.AccountID, &t.Type, &t.Title, &t.Description,
			&t.Amount, &date, &tags)
		if err != nil {
			return nil, err
		}
		t.Tags = strings.Split(tags, ",")
		transactions = append(transactions, t)
	}
	err = rows.Err()
	return transactions, err
}

func (db *Rules) Create(r *waukeen.Rule) error {
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

func (db *Rules) FindAll(acc string) ([]waukeen.Rule, error) {
	var rules []waukeen.Rule

	rows, err := db.Query(`SELECT id, account_id, type, match, result from rules`)
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
