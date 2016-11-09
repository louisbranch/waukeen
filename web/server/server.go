package server

import (
	"fmt"
	"net/http"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/web"
)

type Server struct {
	DB                 waukeen.Database
	Template           web.Template
	StatementsImporter waukeen.StatementsImporter
	RulesImporter      waukeen.RulesImporter
	Transformer        waukeen.TransactionTransformer
	BudgetCalculator   waukeen.BudgetCalculator
}

func (srv *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/accounts/", srv.accounts)
	mux.HandleFunc("/rules/import", srv.importRules)
	mux.HandleFunc("/rules/new", srv.newRule)
	mux.HandleFunc("/rules/", srv.rules)
	mux.HandleFunc("/statements/new", srv.newStatement)
	mux.HandleFunc("/statements", srv.createStatement)
	mux.HandleFunc("/tags/new", srv.newTag)
	mux.HandleFunc("/tags/", srv.tags)
	mux.HandleFunc("/transactions/", srv.transactions)
	mux.HandleFunc("/", srv.index)

	return mux
}

func (srv *Server) render(w http.ResponseWriter, page web.Page) {
	if page.Layout == "" {
		page.Layout = "layout"
	}

	err := srv.Template.Render(w, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
	}
}

func (srv *Server) renderError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	page := web.Page{
		Title:    "500",
		Content:  err,
		Partials: []string{"500"},
	}
	srv.render(w, page)
}

func (srv *Server) renderNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	page := web.Page{
		Title:    "Not Found",
		Partials: []string{"400"},
	}
	srv.render(w, page)
}

func (srv *Server) index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.Redirect(w, r, "/accounts/", http.StatusMovedPermanently)
}
