package search

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/luizbranco/waukeen"
)

func Test_split(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{s: ""},
			want: nil,
		},
		{
			name: "single tag",
			args: args{s: "food"},
			want: []string{"food"},
		},
		{
			name: "multiple tags",
			args: args{s: "food,gifts"},
			want: []string{"food", "gifts"},
		},
		{
			name: "multiple tags with spaces",
			args: args{s: " food , gifts , personal care "},
			want: []string{"food", "gifts", "personal care"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := split(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		v url.Values
		c *http.Cookie
	}
	tests := []struct {
		name string
		args args
		want *Search
	}{
		{
			name: "no form or cookie",
			want: &Search{},
		},
		{
			name: "partial form values",
			args: args{
				v: url.Values{"accounts": []string{"1"}},
			},
			want: &Search{
				Accounts: []string{"1"},
			},
		},
		{
			name: "complete form values",
			args: args{
				v: url.Values{
					"accounts": []string{"1", "2"},
					"types":    []string{"3", "4"},
					"tags":     []string{"food, gift"},
					"start":    []string{"2016-11"},
					"end":      []string{"2016-12"},
				},
			},
			want: &Search{
				Accounts: []string{"1", "2"},
				Types:    []string{"3", "4"},
				Tags:     []string{"food", "gift"},
				Start:    "2016-11",
				End:      "2016-12",
			},
		},
		{
			name: "partial cookie values",
			args: args{
				c: &http.Cookie{
					Name:  "accounts_form",
					Value: "accounts=1",
				},
			},
			want: &Search{
				Accounts: []string{"1"},
			},
		},
		{
			name: "complete cookie values",
			args: args{
				c: &http.Cookie{
					Name:  "accounts_form",
					Value: "accounts=1,2&types=3,4&tags=food, gift&start=2016-11&end=2016-12",
				},
			},
			want: &Search{
				Accounts: []string{"1", "2"},
				Types:    []string{"3", "4"},
				Tags:     []string{"food", "gift"},
				Start:    "2016-11",
				End:      "2016-12",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.Form = tt.args.v
			if tt.args.c != nil {
				r.AddCookie(tt.args.c)
			}

			if got := New(r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearch_DBOptions(t *testing.T) {
	today = func() time.Time {
		return time.Date(2016, 10, 1, 0, 0, 0, 0, time.UTC)
	}

	type fields struct {
		Accounts []string
		Types    []string
		Tags     []string
		Start    string
		End      string
	}
	tests := []struct {
		name   string
		fields fields
		want   waukeen.TransactionsDBOptions
	}{
		{
			name:   "no fields filled",
			fields: fields{},
			want: waukeen.TransactionsDBOptions{
				Types: []waukeen.TransactionType{waukeen.Debit},
				Start: time.Date(2016, 10, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2016, 11, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "some fields filled",
			fields: fields{
				Types: []string{"3", "4"},
				Start: "2016-11",
				End:   "2016-12",
			},
			want: waukeen.TransactionsDBOptions{
				Types: []waukeen.TransactionType{3, 4},
				Start: time.Date(2016, 11, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "all fields filled",
			fields: fields{
				Accounts: []string{"1", "2"},
				Types:    []string{"3", "4"},
				Tags:     []string{"food", "gift"},
				Start:    "2016-11",
				End:      "2016-12",
			},
			want: waukeen.TransactionsDBOptions{
				Accounts: []string{"1", "2"},
				Types:    []waukeen.TransactionType{3, 4},
				Tags:     []string{"food", "gift"},
				Start:    time.Date(2016, 11, 1, 0, 0, 0, 0, time.UTC),
				End:      time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Search{
				Accounts: tt.fields.Accounts,
				Types:    tt.fields.Types,
				Tags:     tt.fields.Tags,
				Start:    tt.fields.Start,
				End:      tt.fields.End,
			}
			if got := f.DBOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search.DBOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearch_empty(t *testing.T) {
	type fields struct {
		Accounts []string
		Types    []string
		Tags     []string
		Start    string
		End      string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "when all fields are not initialized",
			want: true,
		},
		{
			name: "some fields are empty",
			fields: fields{
				Accounts: []string{},
				Tags:     []string{},
			},
			want: true,
		},
		{
			name: "some fields are filled",
			fields: fields{
				Tags: []string{"food"},
			},
			want: false,
		},
		{
			name: "all fields are filled",
			fields: fields{
				Accounts: []string{"1"},
				Types:    []string{"2"},
				Tags:     []string{"food"},
				Start:    "2016-11",
				End:      "2016-11",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Search{
				Accounts: tt.fields.Accounts,
				Types:    tt.fields.Types,
				Tags:     tt.fields.Tags,
				Start:    tt.fields.Start,
				End:      tt.fields.End,
			}
			if got := f.empty(); got != tt.want {
				t.Errorf("Search.empty() = %v, want %v", got, tt.want)
			}
		})
	}
}
