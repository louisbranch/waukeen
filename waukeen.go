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
	Tags        []string
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

type RulesImporter interface {
	Import(io.Reader) ([]Rule, error)
}

type Statement struct {
	Account      Account
	Transactions []Transaction
}

type StatementsImporter interface {
	Import(io.Reader) ([]Statement, error)
}

type Database interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(id string) error
	FindAccount(number string) (*Account, error)
	FindAccounts() ([]Account, error)

	CreateTransaction(t *Transaction) error
	UpdateTransaction(t *Transaction) error
	DeleteTransaction(id string) error
	FindTransaction(id string) (*Transaction, error)
	FindTransactions(TransactionsDBOptions) ([]Transaction, error)

	CreateRule(*Rule) error
	DeleteRule(id string) error
	FindRules(acc string) ([]Rule, error)

	CreateTag(*Tag) error
	DeleteTag(id string) error
	FindTag(name string) (*Tag, error)
	FindTags(starts string) ([]Tag, error)

	CreateStatement(Statement, TransactionTransformer) error
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
