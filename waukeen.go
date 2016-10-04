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
	Create(t *Transaction) error
	FindAll(acc string) ([]Transaction, error)
}

type RulesDB interface {
	Create(*Rule) error
	FindAll(acc string) ([]Rule, error)
}

type TransactionTransformer interface {
	Transform(*Transaction, Rule)
}

var BootstrapTags = []Rule{
	{Type: ReplaceRule, Match: "toronto", Result: ""},
	{Type: TagRule, Match: "pizza", Result: "restaurants"},
	{Type: TagRule, Match: "burger", Result: "restaurants"},
	{Type: TagRule, Match: "burgers", Result: "restaurants"},
	{Type: TagRule, Match: "restaurant", Result: "restaurants"},
	{Type: TagRule, Match: "taco", Result: "restaurants"},
	{Type: TagRule, Match: "tacos", Result: "restaurants"},
	{Type: TagRule, Match: "burrito", Result: "restaurants"},
	{Type: TagRule, Match: "burritos", Result: "restaurants"},
	{Type: TagRule, Match: "subway", Result: "restaurants"},
	{Type: TagRule, Match: "cuisine", Result: "restaurants"},
	{Type: TagRule, Match: "eatery", Result: "restaurants"},
	{Type: TagRule, Match: "sushi", Result: "restaurants"},
	{Type: TagRule, Match: "chipotle", Result: "restaurants"},
	{Type: TagRule, Match: "loblaws", Result: "groceries"},
	{Type: TagRule, Match: "sobeys", Result: "groceries"},
	{Type: TagRule, Match: "cineplex", Result: "entertainment"},
	{Type: TagRule, Match: "LCBO/RAO", Result: "liquor"},
	{Type: TagRule, Match: "ttc", Result: "transportation"},
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
	return "Other"
}
