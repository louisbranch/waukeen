package web

type Link struct {
	Name   string
	URL    string
	Active bool
}

type Page struct {
	Title   string
	Menu    []Link
	Content interface{}
}
