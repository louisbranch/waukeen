package server

import (
	"net/http"
	"strconv"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/web"
	"github.com/pkg/errors"
)

func (srv *Server) newTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	page := web.Page{
		Title:    "New Tag",
		Partials: []string{"tag"},
	}
	srv.render(w, page)
}

func (srv *Server) tags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		name := r.URL.Path[len("/tags/"):]
		if name == "" {
			tags, err := srv.DB.AllTags()

			if err != nil {
				srv.renderError(w, err)
				return
			}

			page := web.Page{
				Title:      "Tags",
				ActiveMenu: "tags",
				Content:    tags,
				Partials:   []string{"tags"},
			}
			srv.render(w, page)
			return
		}

		t, err := srv.DB.FindTag(name)
		if err != nil {
			srv.renderNotFound(w)
			return
		}
		page := web.Page{
			Title:    "Edit Tag",
			Content:  t,
			Partials: []string{"tag"},
		}
		srv.render(w, page)
	case "POST":
		b := r.FormValue("monthly_budget")
		n, err := strconv.ParseInt(b, 10, 64)

		if err != nil {
			errors.Wrap(err, "invalid monthly budget number")
			srv.renderError(w, err)
			return
		}

		id := r.FormValue("id")

		tag := &waukeen.Tag{
			ID:            id,
			Name:          r.FormValue("name"),
			MonthlyBudget: n,
		}

		if id != "" {
			err = srv.DB.UpdateTag(tag)
		} else {
			err = srv.DB.CreateTag(tag)
		}

		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, "/tags/", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
