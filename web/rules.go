package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"

	"github.com/luizbranco/waukeen"
)

func (srv *Server) newRule(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	srv.render(w, nil, "new_rule")
}

func (srv *Server) rules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rules, err := srv.DB.FindRules("")

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

		err = srv.DB.CreateRule(rule)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/rules", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) importRules(w http.ResponseWriter, r *http.Request) {
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
			err := srv.DB.CreateRule(&r)

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
