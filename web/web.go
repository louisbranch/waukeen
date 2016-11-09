package web

import "io"

type Link struct {
	Name   string
	URL    string
	Active bool
}

type Page struct {
	Title    string
	Layout   string
	Partials []string
	Menu     []Link
	Content  interface{}
}

type Template interface {
	Render(w io.Writer, page Page) error
}
