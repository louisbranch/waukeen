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
)

const (
	OtherTransaction TransactionType = iota
	Credit
	Debit
)

const (
	UnknownRule RuleType = iota
	ReplaceRule
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
	Alias       string
	Description string
	Amount      int64
	Date        time.Time
}

type Tag struct {
	ID   string
	Name string
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

type Database interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	FindAccount(number string) (*Account, error)
	FindAccounts() ([]Account, error)

	CreateTransaction(t *Transaction) error
	FindTransactions(TransactionsDBOptions) ([]Transaction, error)

	CreateRule(*Rule) error
	FindRules(acc string) ([]Rule, error)
}

type TransactionsDBOptions struct {
	Accounts []string
	Types    []TransactionType
	Start    time.Time
	End      time.Time
	Tags     []string
}

type TransactionTransformer interface {
	Transform(*Transaction, Rule)
}

type Template interface {
	Render(w io.Writer, data interface{}, paths ...string) error
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

func (t TransactionType) String() string {
	switch t {
	case Debit:
		return "Debit"
	case Credit:
		return "Credit"
	}
	return "Other"
}

func (t RuleType) String() string {
	switch t {
	case ReplaceRule:
		return "Replace"
	case TagRule:
		return "Tagging"
	}
	return "Unknown"
}

func (t *RuleType) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"replace"`:
		*t = ReplaceRule
	case `"tag"`:
		*t = TagRule
	default:
		*t = UnknownRule
	}
	return nil
}
