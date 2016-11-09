package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/web"
)

func (srv *Server) transactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := r.URL.Path[len("/transactions/"):]
		if id == "" {
			srv.renderNotFound(w)
			return
		}
		tr, err := srv.DB.FindTransaction(id)
		if err != nil {
			srv.renderError(w, err)
			return
		}
		if tr.Alias == "" {
			tr.Alias = tr.Title
		}
		page := web.Page{
			Title:    "Transaction",
			Content:  tr,
			Partials: []string{"transaction"},
		}
		srv.render(w, page)
	case "POST":
		id := r.FormValue("id")
		if id == "" {
			srv.renderNotFound(w)
			return
		}
		tr, err := srv.DB.FindTransaction(id)
		if err != nil {
			srv.renderNotFound(w)
			return
		}
		tr.Alias = r.FormValue("alias")
		tr.Description = r.FormValue("description")

		ttype := r.FormValue("transaction_type")
		if ttype != "" {
			i, err := strconv.Atoi(ttype)
			if err == nil {
				tr.Type = waukeen.TransactionType(i)
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
		tr.Tags = tags

		amount := r.FormValue("amount")
		if amount != "" {
			i, err := strconv.ParseInt(amount, 10, 64)
			if err == nil {
				tr.Amount = i
			}
		}

		date := r.FormValue("date")
		if date != "" {
			d, err := time.Parse("2006-01-02", date)
			if err == nil {
				tr.Date = d
			}
		}

		err = srv.DB.UpdateTransaction(tr)

		if err != nil {
			srv.renderError(w, err)
			return
		}

		http.Redirect(w, r, "/accounts", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
