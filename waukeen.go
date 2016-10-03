package waukeen

import (
	"io"
	"time"
)

type AccountType int64

const (
	OtherAccount AccountType = iota
	Checking
	Savings
	CreditCard
)

type Account struct {
	ID       string
	Number   string
	Alias    string
	Type     AccountType
	Currency string
	Balance  int64
}

func (t AccountType) String() string {
	switch t {
	case Checking:
		return "Checking"
	case Savings:
		return "Savings"
	case CreditCard:
		return "Credit Card"
	}
	return "Other"
}

type TransactionType int

const (
	OtherTransaction TransactionType = iota
	Credit
	Debit
)

type Transaction struct {
	ID          string
	FITID       string
	Type        TransactionType
	Name        string
	Description string
	Amount      int64
	Tags        []string
	Date        time.Time
}

type Statement struct {
	Account      Account
	Transactions []Transaction
}

type StatementImporter interface {
	Import(io.Reader) ([]Statement, error)
}

type AccountDB interface {
	Create(*Account) error
	Update(*Account) error
	FindAll() ([]Account, error)
	Find(number string) (*Account, error)
}

type TransactionDB interface {
	FindAll(end time.Time) ([]Transaction, error)
	Create(*Transaction) error
	Update(*Transaction) error
	Delete(*Transaction) error
}
