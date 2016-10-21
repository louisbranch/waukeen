package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/luizbranco/waukeen/html"
	"github.com/luizbranco/waukeen/sqlite"
	"github.com/luizbranco/waukeen/transformer"
	"github.com/luizbranco/waukeen/web"
	"github.com/luizbranco/waukeen/xml"
)

func main() {
	db, err := sqlite.New("waukeen.db")

	if err != nil {
		log.Fatal(err)
	}

	srv := &web.Server{
		DB:          db,
		Template:    html.New("html/templates"),
		Statement:   &xml.XML{},
		Transformer: transformer.Text{},
	}
	mux := srv.NewServeMux()

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
