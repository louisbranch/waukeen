package xml

import (
	"io"

	"github.com/luizbranco/ofx"
	"github.com/luizbranco/waukeen"
)

type Statement struct{}

func (Statement) Import(in io.Reader) ([]waukeen.Statement, error) {
	result, err := ofx.Parse(in)

	if err != nil {
		return nil, err
	}

	var stmts []waukeen.Statement

	for _, r := range result.Banking.BankingResponse {
		res := r.BankStatementResponse
		stmt := waukeen.Statement{
			Account: newBankAccount(res),
		}

		for _, t := range res.BankTransactionsList.Transactions {
			stmt.Transactions = append(stmt.Transactions, newTransaction(t))
		}

		stmts = append(stmts, stmt)
	}

	for _, r := range result.CreditCard.CreditCardResponse {
		res := r.CreditCardStatementResponse
		stmt := waukeen.Statement{
			Account: newCreditCard(res),
		}

		for _, t := range res.BankTransactionsList.Transactions {
			stmt.Transactions = append(stmt.Transactions, newTransaction(t))
		}

		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func newCreditCard(res ofx.CreditCardStatementResponse) waukeen.Account {
	acc := waukeen.Account{
		Type:     waukeen.CreditCard,
		Number:   res.CreditCardAccount.ID,
		Currency: string(res.CurrencyDefault),
		Balance:  int64(res.LedgerBalance.Amount * 100),
	}

	return acc
}

func newBankAccount(res ofx.BankStatementResponse) waukeen.Account {
	acc := waukeen.Account{
		Number:   res.BankingAccount.ID,
		Currency: string(res.CurrencyDefault),
		Balance:  int64(res.LedgerBalance.Amount * 100),
	}

	switch res.BankingAccount.AccountType {
	case ofx.Checking:
		acc.Type = waukeen.Checking
	case ofx.Savings:
		acc.Type = waukeen.Savings
	}

	return acc
}

func newTransaction(res ofx.Transaction) waukeen.Transaction {
	t := waukeen.Transaction{
		FITID:       string(res.FITID),
		Title:       res.Name,
		Description: res.Memo,
		Amount:      int64(res.Amount * 100),
		Date:        res.DatePosted.Time(),
	}

	switch res.TransactionType {
	case ofx.Credit:
		t.Type = waukeen.Credit
	case ofx.Debit:
		t.Type = waukeen.Debit
	}

	return t
}
