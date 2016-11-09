package search

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/luizbranco/waukeen"
)

var today func() time.Time

type Search struct {
	Accounts []string
	Types    []string
	Tags     []string
	Start    string
	End      string
}

func New(r *http.Request) *Search {
	f := &Search{}
	err := r.ParseForm()

	if err == nil {
		f.Accounts = r.Form["accounts"]
		f.Types = r.Form["types"]
		f.Tags = split(r.FormValue("tags"))
		f.Start = r.FormValue("start")
		f.End = r.FormValue("end")

		if !f.empty() {
			return f
		}
	}

	cookie, err := r.Cookie("accounts_form")
	if err != nil {
		return f
	}

	v, err := url.ParseQuery(cookie.Value)

	if err != nil {
		return f
	}

	f.Accounts = split(v.Get("accounts"))
	f.Types = split(v.Get("types"))
	f.Tags = split(v.Get("tags"))
	f.Start = v.Get("start")
	f.End = v.Get("end")

	return f
}

func split(s string) []string {
	var r []string
	vals := strings.Split(s, ",")

	for _, v := range vals {
		v = strings.Trim(v, " ")
		if v != "" {
			r = append(r, v)
		}
	}

	return r
}

func (s *Search) DBOptions() (o waukeen.TransactionsDBOptions) {
	if s.Start != "" {
		t, err := time.Parse("2006-01", s.Start)
		if err == nil {
			o.Start = t
		}
	}

	if s.End != "" {
		t, err := time.Parse("2006-01", s.End)
		if err == nil {
			o.End = t
		}
	}

	for _, t := range s.Types {
		n, err := strconv.Atoi(t)
		if err == nil {
			o.Types = append(o.Types, waukeen.TransactionType(n))
		}
	}

	if len(o.Types) == 0 {
		o.Types = []waukeen.TransactionType{waukeen.Debit}
	}

	o.Accounts = s.Accounts
	o.Tags = s.Tags

	if today == nil {
		today = time.Now
	}

	now := today()

	var year int
	var month time.Month

	if o.Start.IsZero() {
		o.Start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	}

	if o.End.IsZero() {
		month = now.Month()
		year = now.Year()
	} else {
		month = o.End.Month()
		year = o.End.Year()
	}

	if month == time.December {
		o.End = time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)
	} else {
		o.End = time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	}

	s.Start = o.Start.Format("2006-01")
	s.End = o.End.Format("2006-01")

	return o
}

func (f *Search) Save(w http.ResponseWriter) {
	v := make(url.Values)

	for _, e := range f.Accounts {
		v.Add("accounts", e)
	}

	for _, e := range f.Types {
		v.Add("types", e)
	}

	for _, e := range f.Tags {
		v.Add("tags", e)
	}

	v.Set("start", f.Start)
	v.Set("end", f.End)

	cookie := &http.Cookie{
		Name:     "accounts_form",
		Value:    v.Encode(),
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (f *Search) empty() bool {
	return len(f.Accounts) == 0 &&
		len(f.Types) == 0 &&
		len(f.Tags) == 0 &&
		f.Start == "" &&
		f.End == ""
}
