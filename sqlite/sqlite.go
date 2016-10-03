package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/luizbranco/waukeen"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./waukeen.db")
	if err != nil {
		log.Fatal(err)
	}
}

type AccountDB struct {
	*sql.DB
}

func NewAccountDB() (*AccountDB, error) {
	q := `
		CREATE TABLE IF NOT EXISTS accounts(
			id INTEGER PRIMARY KEY,
			number TEXT NOT NULL,
			alias TEXT,
			type INTEGER NOT NULL,
			currency TEXT,
			balance INTEGER
		);
		`

	_, err := db.Exec(q)

	if err != nil {
		return nil, err
	}
	return &AccountDB{db}, nil
}

func (db *AccountDB) Create(a *waukeen.Account) error {
	q := `INSERT into accounts (number, alias, type, currency, balance) values
(?, ?, ?, ?, ?);`

	res, err := db.Exec(q, a.Number, a.Alias, a.Type, a.Currency, a.Balance)

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

func (db *AccountDB) FindAll() ([]waukeen.Account, error) {
	var accounts []waukeen.Account

	rows, err := db.Query("SELECT id, number, alias, type, currency, balance FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		a := waukeen.Account{}
		err = rows.Scan(&a.ID, &a.Number, &a.Alias, &a.Type, &a.Currency, &a.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	err = rows.Err()
	return accounts, err
}

func (db *AccountDB) Find(number string) (*waukeen.Account, error) {
	q := "SELECT id, number, alias, type, currency, balance FROM accounts where number = ?"

	a := &waukeen.Account{}

	err := db.QueryRow(q, number).Scan(&a.ID, &a.Number, &a.Alias, &a.Type,
		&a.Currency, &a.Balance)

	if err != nil {
		return nil, fmt.Errorf("error finding account: %s", err)
	}

	return a, nil
}

func (db *AccountDB) Update(a *waukeen.Account) error {
	_, err := db.Exec(`
	UPDATE accounts SET number=?, alias=?, type=?, currency=?, balance=?
	where id = ?`, a.Number, a.Alias, a.Type, a.Currency, a.Balance, a.ID)
	return err
}
