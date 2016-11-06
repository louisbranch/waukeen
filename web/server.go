package web

import (
	"fmt"
	"net/http"

	"github.com/luizbranco/waukeen"
)

type Server struct {
	DB                 waukeen.Database
	Template           waukeen.Template
	StatementsImporter waukeen.StatementsImporter
	RulesImporter      waukeen.RulesImporter
	Transformer        waukeen.TransactionTransformer
	BudgetCalculator   waukeen.BudgetCalculator
}

func (srv *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

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

func (srv *Server) render(w http.ResponseWriter, data interface{}, path ...string) {
	err := srv.Template.Render(w, data, path...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
	}
}

func (srv *Server) renderError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	srv.render(w, err, "500")
}

func (srv *Server) index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	srv.render(w, nil, "index")
}
