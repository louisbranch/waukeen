package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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
			account_id INTEGER NOT NULL,
			fitid TEXT NOT NULL CHECK(fitid <> ''),
			type INTEGER NOT NULL,
			title TEXT NOT NULL CHECK(title <> ''),
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
			name TEXT NOT NULL UNIQUE CHECK(name <> '')
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS transaction_tags(
			id INTEGER PRIMARY KEY,
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
			match TEXT NOT NULL CHECK(match <> ''),
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

func (db *DB) DeleteAccount(id string) error {
	return errors.New("not implemented")
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
	q := `INSERT into transactions
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

	for _, name := range t.Tags {
		tag, err := db.FindTag(name)
		if err != nil {
			tag = &waukeen.Tag{Name: name}
			err = db.CreateTag(tag)
		}
		if err != nil {
			return fmt.Errorf("error creating transaction tag %s:", name, err)
		}
		q := `INSERT into transaction_tags (transaction_id, tag_id) values (?, ?)`
		_, err = db.Exec(q, t.ID, tag.ID)
		if err != nil {
			return fmt.Errorf("error creating transaction tag relation %s:", name, err)
		}
	}

	return nil
}

func (db *DB) UpdateTransaction(t *waukeen.Transaction) error {
	q := fmt.Sprintf(`delete from transaction_tags where id in (select
		transaction_tags.id from transaction_tags INNER JOIN tags on tags.id =
		transaction_tags.tag_id where transaction_tags.transaction_id = ? AND
		tags.name NOT IN (%s))`, toInCodition(t.Tags))

	_, err := db.Exec(q, t.ID)

	if err != nil {
		return err
	}

	for _, name := range t.Tags {
		tag, err := db.FindTag(name)
		if err != nil {
			tag = &waukeen.Tag{Name: name}
			err = db.CreateTag(tag)
		}
		if err != nil {
			return fmt.Errorf("error updating transaction tag %s:", name, err)
		}
		q := `INSERT OR IGNORE into transaction_tags (transaction_id, tag_id) values (?, ?)`
		_, err = db.Exec(q, t.ID, tag.ID)
		if err != nil {
			return fmt.Errorf("error updating transaction tag relation %s:", name, err)
		}
	}

	q = `UPDATE transactions SET account_id = ?, fitid = ?, type = ?, title = ?,
	alias = ?, description = ?, amount = ?, date = ? WHERE id = ?`

	_, err = db.Exec(q, t.AccountID, t.FITID, t.Type, t.Title, t.Alias,
		t.Description, t.Amount, t.Date, t.ID)

	if err != nil {
		return fmt.Errorf("error updating transaction: %s", err)
	}

	return nil
}

func (db *DB) DeleteTransaction(id string) error {
	return errors.New("not implemented")
}

func (db *DB) DeleteRule(id string) error {
	return errors.New("not implemented")
}

func (db *DB) DeleteTag(id string) error {
	return errors.New("not implemented")
}

func (db *DB) FindTransactions(opts waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error) {
	var transactions []waukeen.Transaction
	var query string
	var clauses []string

	if len(opts.Tags) > 0 {
		query = `SELECT transactions.id, transactions.account_id,
		transactions.fitid, transactions.type, transactions.title,
		transactions.alias, transactions.description, transactions.amount,
		transactions.date FROM transactions JOIN transaction_tags ON
		transactions.id = transaction_tags.transaction_id JOIN tags ON tags.id
		= transaction_tags.tag_id where `
		tags := toInCodition(opts.Tags)
		clauses = append(clauses, fmt.Sprintf("tags.name IN (%s)", tags))
	} else {
		query = `SELECT id, account_id, fitid, type, title, alias, description, amount, date
	FROM transactions WHERE `
	}

	if len(opts.Accounts) > 0 {
		clause := "transactions.account_id IN (" + strings.Join(opts.Accounts, " ,") + ")"
		clauses = append(clauses, clause)
	}

	if len(opts.Types) > 0 {
		var types []string

		for _, t := range opts.Types {
			types = append(types, strconv.Itoa(int(t)))
		}

		clause := "transactions.type IN (" + strings.Join(types, " ,") + ")"
		clauses = append(clauses, clause)
	}

	if !opts.Start.IsZero() {
		clauses = append(clauses, "transactions.date >= "+opts.Start.Format("'2006-01-02'"))
	}

	if !opts.End.IsZero() {
		end := opts.End
		end = end.Add(time.Hour * 24)
		clauses = append(clauses, "transactions.date < "+end.Format("'2006-01-02'"))
	}

	query += strings.Join(clauses, " AND ")

	if len(opts.Tags) > 0 {
		query += "GROUP BY transactions.id"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := waukeen.Transaction{}
		err = rows.Scan(&t.ID, &t.AccountID, &t.FITID, &t.Type, &t.Title, &t.Alias,
			&t.Description, &t.Amount, &t.Date)
		if err != nil {
			return nil, err
		}

		tags, err := db.findTags(t.ID)
		if err != nil {
			return nil, err
		}
		t.Tags = tags

		transactions = append(transactions, t)
	}
	err = rows.Err()
	return transactions, err
}

func (db *DB) findTags(transaction string) ([]string, error) {
	q := `SELECT DISTINCT tags.name FROM transaction_tags JOIN tags on
	transaction_tags.tag_id = tags.id WHERE transaction_tags.transaction_id = ?`

	var tags []string

	rows, err := db.Query(q, transaction)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, name)
	}
	err = rows.Err()
	sort.Strings(tags)

	return tags, err
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
	where account_id = ? OR account_id = ''`

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

func (db *DB) CreateTag(t *waukeen.Tag) error {
	q := `INSERT into tags (name) values (?)`

	res, err := db.Exec(q, t.Name, t.Name)

	if err != nil {
		return fmt.Errorf("error creating tag: %s", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return fmt.Errorf("error retrieving last tag id: %s", err)
	}

	t.ID = strconv.FormatInt(id, 10)

	return nil
}

func (db *DB) FindTag(name string) (*waukeen.Tag, error) {
	q := "SELECT id, name FROM tags where name = ?"

	t := &waukeen.Tag{}

	err := db.QueryRow(q, name).Scan(&t.ID, &t.Name)

	if err != nil {
		return nil, fmt.Errorf("error finding tag: %s", err)
	}

	return t, nil
}

func (db *DB) FindTags(starts string) ([]waukeen.Tag, error) {
	var tags []waukeen.Tag

	starts += "%"
	rows, err := db.Query("SELECT id, name FROM tags where name LIKE ?", starts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := waukeen.Tag{}
		err = rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	err = rows.Err()
	return tags, err
}

func toInCodition(args []string) string {
	return "'" + strings.Join(args, "', '") + "'"
}
