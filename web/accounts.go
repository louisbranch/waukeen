package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/luizbranco/waukeen"
)

func getTransactionForm(r *http.Request) waukeen.TransactionsDBOptions {
	opt := waukeen.TransactionsDBOptions{}

	start := r.FormValue("start")
	if start != "" {
		t, err := time.Parse("2006-01-02", start)
		if err == nil {
			opt.Start = t
		}
	}

	end := r.FormValue("end")
	if end != "" {
		t, err := time.Parse("2006-01-02", end)
		if err == nil {
			opt.End = t
		}
	}

	number := r.FormValue("account")
	if number != "" {
		opt.Accounts = append(opt.Accounts, number)
	}

	ttype := r.FormValue("transaction_type")
	if ttype != "" {
		i, err := strconv.Atoi(ttype)
		if err == nil {
			opt.Types = []waukeen.TransactionType{waukeen.TransactionType(i)}
		}
	}

	var tags []string
	vals := strings.Split(r.FormValue("tags"), ",")
	for _, t := range vals {
		tag := strings.Trim(t, " ")
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	if len(tags) > 0 {
		opt.Tags = tags
	}

	return opt
}

func (srv *Server) accounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	opt := getTransactionForm(r)

	accounts, err := srv.DB.FindAccounts(opt.Accounts...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	transactions, err := srv.DB.FindTransactions(opt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	var total int64
	for _, t := range transactions {
		total += t.Amount
	}

	content := struct {
		Form         waukeen.TransactionsDBOptions
		Accounts     []waukeen.Account
		Transactions []waukeen.Transaction
		Total        int64
	}{
		Form:         opt,
		Accounts:     accounts,
		Transactions: transactions,
		Total:        total,
	}

	srv.render(w, content, "accounts")
}
