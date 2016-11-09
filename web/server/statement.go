package server

import "net/http"

func (srv *Server) newStatement(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	srv.render(w, nil, "statement")
}

func (srv *Server) createStatement(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	file, _, err := r.FormFile("statement")
	if err != nil {
		srv.renderError(w, err)
		return
	}

	list, err := srv.StatementsImporter.Import(file)
	if err != nil {
		srv.renderError(w, err)
		return
	}

	if len(list) == 0 {
		srv.renderError(w, err)
		return
	}

	for _, item := range list {
		err = srv.DB.CreateStatement(item, srv.Transformer)
		if err != nil {
			srv.renderError(w, err)
			return
		}
	}

	http.Redirect(w, r, "/accounts", http.StatusFound)
}
