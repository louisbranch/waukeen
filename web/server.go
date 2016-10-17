package web

import (
	"encoding/json"
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

type Server struct {
	Statement    waukeen.StatementImporter
	Accounts     waukeen.AccountsDB
	Transactions waukeen.TransactionsDB
	Rules        waukeen.RulesDB
	Transformer  waukeen.TransactionTransformer
}

type AccountContent struct {
	Account      *waukeen.Account
	Total        int64
	Transactions []waukeen.Transaction
}

func (srv *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/accounts", srv.accounts)
	mux.HandleFunc("/statements", srv.statementCreate)
	mux.HandleFunc("/statements/new", srv.statementNew)
	mux.HandleFunc("/rules/batch", srv.rulesBatch)
	mux.HandleFunc("/rules/new", srv.rulesNew)
	mux.HandleFunc("/rules", srv.rules)
	mux.HandleFunc("/", srv.index)

	return mux
}

func (srv *Server) statementNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	p := path.Join("web", "templates", "statement.html")
	t, err := template.ParseFiles(p)
	if err == nil {
		err = t.Execute(w, nil)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (srv *Server) statementCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	file, _, err := r.FormFile("statement")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	list, err := srv.Statement.Import(file)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	for _, stmt := range list {
		number := stmt.Account.Number

		acc, err := srv.Accounts.Find(number)

		if err == nil {
			acc.Balance = stmt.Account.Balance
			err = srv.Accounts.Update(acc)
		} else {
			acc = &stmt.Account
			err = srv.Accounts.Create(acc)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, err)
				return
			}

		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		rules, err := srv.Rules.FindAll(acc.ID)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		for _, tn := range stmt.Transactions {
			t := &tn
			t.AccountID = acc.ID
			for _, r := range rules {
				srv.Transformer.Transform(t, r)
			}
			err := srv.Transactions.Create(t)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, err)
				return
			}
		}

	}
	http.Redirect(w, r, "/accounts", http.StatusFound)
}

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
		acc, err := srv.Accounts.Find(number)
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
		accs, err = srv.Accounts.FindAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

	}

	for _, a := range accs {
		opts.Accounts = []string{a.ID}
		transactions, err := srv.Transactions.FindAll(opts)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		var total int64

		for _, t := range transactions {
			total += t.Amount
		}

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

func (srv *Server) rulesNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	p := path.Join("web", "templates", "new_rule.html")
	t, err := template.ParseFiles(p)
	if err == nil {
		err = t.Execute(w, nil)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (srv *Server) rules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rules, err := srv.Rules.FindAll("")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		p := path.Join("web", "templates", "rules.html")
		t, err := template.ParseFiles(p)
		if err == nil {
			err = t.Execute(w, rules)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "POST":
		t := r.FormValue("type")
		n, err := strconv.Atoi(t)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rule := &waukeen.Rule{
			AccountID: r.FormValue("account"),
			Type:      waukeen.RuleType(n),
			Match:     r.FormValue("match"),
			Result:    r.FormValue("result"),
		}

		if rule.Match == "" {
			http.Redirect(w, r, "/rules/new", http.StatusFound)
			return
		}

		err = srv.Rules.Create(rule)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/rules", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) rulesBatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		p := path.Join("web", "templates", "batch_rules.html")
		t, err := template.ParseFiles(p)
		if err == nil {
			err = t.Execute(w, nil)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "POST":
		file, _, err := r.FormFile("rules")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		var rules []waukeen.Rule

		dec := json.NewDecoder(file)

		err = dec.Decode(&rules)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		for _, r := range rules {
			err := srv.Rules.Create(&r)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, err)
				return
			}

		}

		http.Redirect(w, r, "/rules", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	p := path.Join("web", "templates", "index.html")
	t, err := template.ParseFiles(p)
	if err == nil {
		err = t.Execute(w, nil)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
