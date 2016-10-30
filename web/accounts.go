package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/luizbranco/waukeen"
)

type accountForm struct {
	Accounts []string
	Types    []string
	Tags     []string
	Start    string
	End      string
}

func newAccountForm(opt waukeen.TransactionsDBOptions) accountForm {
	acc := accountForm{
		Accounts: opt.Accounts,
		Tags:     opt.Tags,
	}

	acc.Start = opt.Start.Format("2006-01-02")
	acc.End = opt.End.Format("2006-01-02")

	for _, t := range opt.Types {
		acc.Types = append(acc.Types, strconv.Itoa(int(t)))
	}

	return acc
}

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

	opt.Accounts = r.Form["account"]

	ttype, ok := r.Form["transaction_type"]
	if ok {
		opt.Types = make([]waukeen.TransactionType, len(ttype))
		for i, t := range ttype {
			n, err := strconv.Atoi(t)
			if err == nil {
				opt.Types[i] = waukeen.TransactionType(n)
			}
		}
	} else {
		opt.Types = []waukeen.TransactionType{waukeen.Debit}
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

	now := time.Now()
	year := now.Year()
	month := now.Month()
	startT := opt.Start
	endT := opt.End

	if startT.IsZero() {
		startT = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	}

	if endT.IsZero() {
		if month == time.December {
			endT = time.Date(year, month, 31, 0, 0, 0, 0, time.UTC)
		} else {
			endT = time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC).Add(-24 * time.Hour)
		}
	}

	opt.Start = startT
	opt.End = endT

	return opt
}

func (srv *Server) accounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	opt := getTransactionForm(r)

	accounts, err := srv.DB.FindAccounts()
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

	ids := make([]string, len(accounts))

	for i, acc := range accounts {
		ids[i] = acc.ID
	}

	opt.Accounts = ids

	content := struct {
		Form         accountForm
		Accounts     []waukeen.Account
		Transactions []waukeen.Transaction
		Total        int64
	}{
		Form:         newAccountForm(opt),
		Accounts:     accounts,
		Transactions: transactions,
		Total:        total,
	}

	srv.render(w, content, "accounts")
}
