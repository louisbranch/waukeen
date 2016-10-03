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
		switch r.Method {

		case "GET":
			accs, err := srv.DB.FindAll()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err)
				return
			}

			for _, a := range accs {
				fmt.Fprintln(w, a)
			}

		case "POST":
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

				acc, err := srv.DB.Find(number)

				if err == nil {
					acc.Balance = stmt.Account.Balance
					err = srv.DB.Update(acc)
				} else {
					acc = &stmt.Account
					err = srv.DB.Create(acc)
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

				fmt.Fprintf(w, "Account %s (%s): %s %.2f\n", acc.Number, acc.Type,
					acc.Currency, float64(acc.Balance)/100)
				for _, t := range stmt.Transactions {
					fmt.Fprintf(w, "\t%s: %.2f\n", t.Name, float64(t.Amount)/100)
				}
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
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
