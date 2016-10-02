package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/luizbranco/waukeen"
)

type Server struct {
	Statement waukeen.StatementImporter
	DB        waukeen.AccountDB
}

func (srv *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/statements", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		file, _, err := r.FormFile("statement")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error: %s", err)
			return
		}

		list, err := srv.Statement.Import(file)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error: %s", err)
			return
		}

		for _, stmt := range list {
			acc := stmt.Account
			fmt.Fprintf(w, "Account %s (%s): %s %.2f\n", acc.Number, acc.Type,
				acc.Currency, float64(acc.Balance)/100)
			for _, t := range stmt.Transactions {
				fmt.Fprintf(w, "\t%s: %.2f\n", t.Name, float64(t.Amount)/100)
			}
		}

	})

	mux.HandleFunc("/statements/new", func(w http.ResponseWriter, r *http.Request) {
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
	})

	return mux
}
