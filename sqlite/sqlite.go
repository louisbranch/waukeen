package sqlite

import (
	"database/sql"
	"log"

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

func NewAccountDB() *AccountDB {
	var name string
	q := "SELECT name FROM sqlite_master WHERE type='table' AND name='?';"

	err := db.QueryRow(q).Scan(&name)

	if err != nil {
		//TODO: create table
	}
	return &AccountDB{db}
}

func (db *AccountDB) All() []waukeen.Account {

	return nil
}

func (db *AccountDB) Find(number string) (*waukeen.Account, error) {

	return nil, nil
}
