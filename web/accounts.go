package web

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/luizbranco/waukeen"
)

func (srv *Server) accounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	opts := waukeen.TransactionsDBOptions{}

	start := r.FormValue("start")
	if start != "" {
		t, err := time.Parse("2006-01-02", start)
		if err == nil {
			opts.Start = t
		}
	}

	end := r.FormValue("end")
	if end != "" {
		t, err := time.Parse("2006-01-02", end)
		if err == nil {
			opts.End = t
		}
	}

	var accs []waukeen.Account
	var err error
	number := r.FormValue("account")

	if number != "" {
		acc, err := srv.DB.FindAccount(number)
		if err == nil {
			accs = append(accs, *acc)
		}
	}

	ttype := r.FormValue("transaction_type")
	if ttype != "" {
		i, err := strconv.Atoi(ttype)
		if err == nil {
			opts.Types = []waukeen.TransactionType{waukeen.TransactionType(i)}
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
	opts.Tags = tags

	type form struct {
		Account string
		Start   string
		End     string
		Type    string
		Tags    []string
	}

	content := struct {
		Form           form
		AccountContent []AccountContent
	}{
		Form: form{Account: number, Start: start, End: end, Type: ttype, Tags: tags},
	}

	if len(accs) == 0 {
		accs, err = srv.DB.FindAccounts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

	}

	for _, a := range accs {
		opts.Accounts = []string{a.ID}
		transactions, err := srv.DB.FindTransactions(opts)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		var total int64

		c := AccountContent{
			Account:      &a,
			Total:        total,
			Transactions: transactions,
		}

		content.AccountContent = append(content.AccountContent, c)
	}

	fns := template.FuncMap{"currency": func(amount int64) string {
		return fmt.Sprintf("$%.2f", math.Abs(float64(amount))/100)
	}}

	t := template.New("").Funcs(fns)

	p := path.Join("web", "templates", "accounts.html")
	t, err = t.ParseFiles(p)
	if err == nil {
		t = t.Funcs(fns)
		err = t.ExecuteTemplate(w, "accounts.html", content)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	}
}
