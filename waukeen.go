package waukeen

import (
	"io"
	"time"
)

type AccountType int
type TransactionType int
type RuleType int

const (
	OtherAccount AccountType = iota
	Checking
	Savings
	CreditCard

	OtherTransaction TransactionType = iota
	Credit
	Debit

	ReplaceRule RuleType = iota
	TagRule
)

type Account struct {
	ID       string
	Number   string
	Name     string
	Type     AccountType
	Currency string
	Balance  int64
}

type Transaction struct {
	ID          string
	AccountID   string
	FITID       string
	Type        TransactionType
	Title       string
	Description string
	Amount      int64
	Tags        []string
	Date        time.Time
}

type Rule struct {
	ID        string
	AccountID string
	Type      RuleType
	Match     string
	Result    string
}

type Statement struct {
	Account      Account
	Transactions []Transaction
}

type StatementImporter interface {
	Import(io.Reader) ([]Statement, error)
}

type AccountsDB interface {
	Create(*Account) error
	Update(*Account) error
	FindAll() ([]Account, error)
	Find(number string) (*Account, error)
}

type TransactionsDB interface {
	Create(acc string, t *Transaction) error
	FindAll(acc string) ([]Transaction, error)
	//Update(*Transaction) error
	//Find(FITID string) (*Transaction, error)
	//FindAll(end time.Time) ([]Transaction, error)
	//Delete(*Transaction) error
}

type RulesDB interface {
	Create(*Rule) error
	FindAll(acc string) ([]Rule, error)
}

type TransactionTransformer interface {
	Transform(*Transaction, Rule)
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
