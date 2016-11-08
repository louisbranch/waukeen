package web

import (
	"net/http"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/web/search"
)

func (srv *Server) accounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	form := search.New(r)
	opt := form.DBOptions()

	accs, err := srv.DB.FindAccounts()
	if err != nil {
		srv.renderError(w, err)
		return
	}

	transactions, err := srv.DB.FindTransactions(opt)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	tags, err := srv.DB.AllTags()
	if err != nil {
		srv.renderError(w, err)
		return
	}

	var total int64
	for _, t := range transactions {
		total += t.Amount
	}

	ids := make([]string, len(accs))

	for i, acc := range accs {
		ids[i] = acc.ID
	}

	form.Accounts = ids
	opt.Accounts = ids

	months := monthSpam(opt)
	budgets := srv.BudgetCalculator.Calculate(months, transactions, tags)

	content := struct {
		Form         *search.Search
		Accounts     []waukeen.Account
		Transactions []waukeen.Transaction
		Total        int64
		Budgets      []waukeen.Budget
	}{
		Form:         form,
		Accounts:     accs,
		Transactions: transactions,
		Total:        total,
		Budgets:      budgets,
	}

	form.Save(w)

	srv.render(w, content, "accounts")
}

func monthSpam(opt waukeen.TransactionsDBOptions) int {
	years := opt.End.Year() - opt.Start.Year()
	months := (int(opt.End.Month()) + (years * 12)) - int(opt.Start.Month())
	return months + 1
}
