package server

import (
	"net/http"
	"strconv"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/web"
)

func (srv *Server) newRule(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	page := web.Page{
		Title:    "New Rule",
		Partials: []string{"new_rule"},
	}

	srv.render(w, page)
}

func (srv *Server) rules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rules, err := srv.DB.FindRules()

		if err != nil {
			srv.renderError(w, err)
			return
		}

		page := web.Page{
			Title:      "Rules",
			ActiveMenu: "rules",
			Content:    rules,
			Partials:   []string{"rules"},
		}

		srv.render(w, page)
	case "POST":
		t := r.FormValue("type")
		n, err := strconv.Atoi(t)

		if err != nil {
			srv.renderError(w, err)
			return
		}

		rule := &waukeen.Rule{
			Type:   waukeen.RuleType(n),
			Match:  r.FormValue("match"),
			Result: r.FormValue("result"),
		}

		err = srv.DB.CreateRule(rule)

		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, "/rules/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) importRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		page := web.Page{
			Title:    "Import Rules",
			Partials: []string{"import_rules"},
		}
		srv.render(w, page)
	case "POST":
		file, _, err := r.FormFile("rules")
		if err != nil {
			srv.renderError(w, err)
			return
		}

		rules, err := srv.RulesImporter.Import(file)

		if err != nil {
			srv.renderError(w, err)
			return
		}

		for _, r := range rules {
			err := srv.DB.CreateRule(&r)

			if err != nil {
				srv.renderError(w, err)
				return
			}

		}

		http.Redirect(w, r, "/rules", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
