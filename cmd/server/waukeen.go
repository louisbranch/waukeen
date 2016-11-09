package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/luizbranco/waukeen/calc"
	"github.com/luizbranco/waukeen/json"
	"github.com/luizbranco/waukeen/sqlite"
	"github.com/luizbranco/waukeen/transformer"
	"github.com/luizbranco/waukeen/web/html"
	"github.com/luizbranco/waukeen/web/server"
	"github.com/luizbranco/waukeen/xml"
)

func main() {
	db, err := sqlite.New("waukeen.db")

	if err != nil {
		log.Fatal(err)
	}

	srv := &server.Server{
		DB:                 db,
		Template:           html.New("web/templates"),
		StatementsImporter: xml.Statement{},
		RulesImporter:      json.Rules{},
		Transformer:        transformer.Text{},
		BudgetCalculator:   calc.Budgeter{},
	}
	mux := srv.NewServeMux()

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
