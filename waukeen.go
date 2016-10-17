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

type TransactionsDBOptions struct {
	Account string
	Start   time.Time
	End     time.Time
	Tags    []string
}

type TransactionsDB interface {
	Create(t *Transaction) error
	FindAll(TransactionsDBOptions) ([]Transaction, error)
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
