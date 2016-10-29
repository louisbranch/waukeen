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

type RulesImporter struct {
	ImportMethod func(io.Reader) ([]waukeen.Rule, error)
}

func (m *RulesImporter) Import(in io.Reader) ([]waukeen.Rule, error) {
	return m.ImportMethod(in)
}

type StatementsImporter struct {
	ImportMethod func(io.Reader) ([]waukeen.Statement, error)
}

func (m *StatementsImporter) Import(in io.Reader) ([]waukeen.Statement, error) {
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
	DeleteAccountMethod func(string) error
	FindAccountMethod   func(number string) (*waukeen.Account, error)
	FindAccountsMethod  func() ([]waukeen.Account, error)

	CreateTransactionMethod func(*waukeen.Transaction) error
	UpdateTransactionMethod func(*waukeen.Transaction) error
	DeleteTransactionMethod func(string) error
	FindTransactionsMethod  func(waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error)
	FindTransactionMethod   func(string) (*waukeen.Transaction, error)

	CreateRuleMethod func(*waukeen.Rule) error
	DeleteRuleMethod func(string) error
	FindRulesMethod  func(acc string) ([]waukeen.Rule, error)

	CreateTagMethod func(*waukeen.Tag) error
	DeleteTagMethod func(string) error
	FindTagMethod   func(name string) (*waukeen.Tag, error)
	FindTagsMethod  func(starts string) ([]waukeen.Tag, error)

	CreateStatementMethod func(waukeen.Statement, waukeen.TransactionTransformer) error
}

func (m *Database) CreateAccount(a *waukeen.Account) error {
	return m.CreateAccountMethod(a)
}

func (m *Database) DeleteAccount(id string) error {
	return m.DeleteAccountMethod(id)
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

func (m *Database) UpdateTransaction(t *waukeen.Transaction) error {
	return m.UpdateTransactionMethod(t)
}

func (m *Database) DeleteTransaction(id string) error {
	return m.DeleteTransactionMethod(id)
}

func (m *Database) FindTransactions(opts waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error) {
	return m.FindTransactionsMethod(opts)
}

func (m *Database) FindTransaction(id string) (*waukeen.Transaction, error) {
	return m.FindTransactionMethod(id)
}

func (m *Database) CreateRule(r *waukeen.Rule) error {
	return m.CreateRuleMethod(r)
}

func (m *Database) DeleteRule(id string) error {
	return m.DeleteRuleMethod(id)
}

func (m *Database) FindRules(acc string) ([]waukeen.Rule, error) {
	return m.FindRulesMethod(acc)
}

func (m *Database) CreateStatement(s waukeen.Statement, t waukeen.TransactionTransformer) error {
	return m.CreateStatementMethod(s, t)
}

func (m *Database) CreateTag(t *waukeen.Tag) error {
	return m.CreateTagMethod(t)
}

func (m *Database) FindTag(name string) (*waukeen.Tag, error) {
	return m.FindTagMethod(name)
}

func (m *Database) FindTags(starts string) ([]waukeen.Tag, error) {
	return m.FindTagsMethod(starts)
}

func (m *Database) DeleteTag(id string) error {
	return m.DeleteTagMethod(id)
}
