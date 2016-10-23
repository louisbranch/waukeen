package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/luizbranco/waukeen"
)

type Server struct {
	DB          waukeen.Database
	Template    waukeen.Template
	Statement   waukeen.StatementImporter
	Transformer waukeen.TransactionTransformer
}

type TagCost struct {
	Name  string
	Count int
	Total int64
}

type TagCosts []TagCost

type AccountContent struct {
	Account      *waukeen.Account
	Total        int64
	Transactions []waukeen.Transaction
	TagCosts     []TagCost
}

func (srv *Server) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/accounts", srv.accounts)
	mux.HandleFunc("/rules/batch", srv.importRules)
	mux.HandleFunc("/rules/new", srv.newRule)
	mux.HandleFunc("/rules", srv.rules)
	mux.HandleFunc("/statements", srv.createStatement)
	mux.HandleFunc("/statements/new", srv.newStatement)
	mux.HandleFunc("/", srv.index)

	return mux
}

func (srv *Server) render(w http.ResponseWriter, data interface{}, path ...string) {
	err := srv.Template.Render(w, data, path...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
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

func (t TagCosts) Len() int {
	return len(t)
}

func (t TagCosts) Less(i, j int) bool {
	if t[i].Name == "others" {
		return false
	}

	if t[j].Name == "others" {
		return true
	}

	return strings.Compare(t[i].Name, t[j].Name) < 0
}

func (t TagCosts) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
