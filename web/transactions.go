package web

import "net/http"

func (srv *Server) transactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := r.URL.Path[len("/transactions/"):]
		if id == "" {
			srv.render(w, nil, "404")
			return
		}

	case "POST":
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
