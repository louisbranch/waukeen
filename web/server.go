package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"

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
	Transactions []waukeen.Transaction
}

func (srv *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/accounts", srv.accountsIndex)
	mux.HandleFunc("/statements", srv.statementCreate)
	mux.HandleFunc("/statements/new", srv.statementNew)
	mux.HandleFunc("/rules/new", srv.rulesNew)
	mux.HandleFunc("/rules", srv.rules)

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

			for _, r := range waukeen.BootstrapTags {
				r.AccountID = acc.ID
				err := srv.Rules.Create(&r)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, err)
					return
				}
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

func (srv *Server) accountsIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	content := struct {
		AccountContent []AccountContent
	}{}

	accs, err := srv.Accounts.FindAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	for _, a := range accs {
		transactions, err := srv.Transactions.FindAll(a.ID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		c := AccountContent{
			Account:      &a,
			Transactions: transactions,
		}

		content.AccountContent = append(content.AccountContent, c)
	}

	p := path.Join("web", "templates", "accounts.html")
	t, err := template.ParseFiles(p)
	if err == nil {
		err = t.Execute(w, content)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
