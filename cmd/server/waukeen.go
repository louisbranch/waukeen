package main

import (
	"log"
	"net/http"

	"github.com/luizbranco/waukeen/sqlite"
	"github.com/luizbranco/waukeen/web"
	"github.com/luizbranco/waukeen/xml"
)

func main() {
	importer := &xml.XML{}
	db := sqlite.NewAccountDB()

	srv := &web.Server{
		Statement: importer,
		DB:        db,
	}
	mux := srv.NewServeMux()

	log.Fatal(http.ListenAndServe(":8080", mux))
}
