package mock

import (
	"io"

	"github.com/luizbranco/waukeen"
)

type Template struct {
	RenderMethod func(w io.Writer, data interface{}, path ...string) error
}

func (m *Template) Render(w io.Writer, data interface{}, path ...string) error {
	return m.RenderMethod(w, data, path...)
}

type StatementImporter struct {
	ImportMethod func(io.Reader) ([]waukeen.Statement, error)
}

func (m *StatementImporter) Import(in io.Reader) ([]waukeen.Statement, error) {
	return m.ImportMethod(in)
}

type TransactionTransformer struct {
	TransformMethod func(*waukeen.Transaction, waukeen.Rule)
}

func (m *TransactionTransformer) Transform(t *waukeen.Transaction, r waukeen.Rule) {
	m.TransformMethod(t, r)
}

type Database struct {
	CreateAccountMethod func(*waukeen.Account) error
	UpdateAccountMethod func(*waukeen.Account) error
	FindAccountMethod   func(number string) (*waukeen.Account, error)
	FindAccountsMethod  func() ([]waukeen.Account, error)

	CreateTransactionMethod func(t *waukeen.Transaction) error
	FindTransactionsMethod  func(waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error)

	CreateRuleMethod func(*waukeen.Rule) error
	FindRulesMethod  func(acc string) ([]waukeen.Rule, error)

	CreateStatementMethod func(waukeen.Statement, waukeen.TransactionTransformer) error
}

func (m *Database) CreateAccount(a *waukeen.Account) error {
	return m.CreateAccountMethod(a)
}

func (m *Database) UpdateAccount(a *waukeen.Account) error {
	return m.UpdateAccountMethod(a)
}

func (m *Database) FindAccount(number string) (*waukeen.Account, error) {
	return m.FindAccountMethod(number)
}

func (m *Database) FindAccounts() ([]waukeen.Account, error) {
	return m.FindAccountsMethod()
}

func (m *Database) CreateTransaction(t *waukeen.Transaction) error {
	return m.CreateTransactionMethod(t)
}

func (m *Database) FindTransactions(opts waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error) {
	return m.FindTransactionsMethod(opts)
}

func (m *Database) CreateRule(r *waukeen.Rule) error {
	return m.CreateRuleMethod(r)
}

func (m *Database) FindRules(acc string) ([]waukeen.Rule, error) {
	return m.FindRulesMethod(acc)
}

func (m *Database) CreateStatement(s waukeen.Statement, t waukeen.TransactionTransformer) error {
	return m.CreateStatementMethod(s, t)
}
