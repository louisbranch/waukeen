package web

import (
	"errors"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/mock"
)

func TestAccounts(t *testing.T) {
	db := &mock.Database{}
	budgeter := &mock.BudgetCalculator{}

	db.FindAccountsMethod = func(ids ...string) ([]waukeen.Account, error) {
		return nil, nil
	}
	db.FindTransactionsMethod = func(waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error) {
		return nil, nil
	}
	db.AllTagsMethod = func() ([]waukeen.Tag, error) {
		return nil, nil
	}

	budgeter.CalculateMethod = func(trs []waukeen.Transaction,
		tags []waukeen.Tag) []waukeen.Budget {
		return nil
	}

	srv := &Server{DB: db, BudgetCalculator: budgeter}

	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/accounts", nil)
		res := serverTest(nil, req)
		code := 405
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Empty accounts and transactions list", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/accounts", nil)
		res := serverTest(srv, req)
		code := 200
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Find accounts DB error", func(t *testing.T) {
		db.FindAccountsMethod = func(ids ...string) ([]waukeen.Account, error) {
			return nil, errors.New("not implemented")
		}
		req := httptest.NewRequest("GET", "/accounts", nil)
		res := serverTest(srv, req)

		code := 500
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Find transactions DB error", func(t *testing.T) {
		db.FindTransactionsMethod = func(waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error) {
			return nil, errors.New("not implemented")
		}
		req := httptest.NewRequest("GET", "/accounts", nil)
		res := serverTest(srv, req)

		code := 500
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Empty Form", func(t *testing.T) {
		db.FindAccountsMethod = func(ids ...string) ([]waukeen.Account, error) {
			return []waukeen.Account{{ID: "1"}}, nil
		}
		db.FindTransactionsMethod = func(got waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error) {
			/* FIXME
			want := waukeen.TransactionsDBOptions{}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("wants %+v, got %+v", want, got)
			}
			*/
			return nil, nil
		}
		req := httptest.NewRequest("GET", "/accounts", nil)
		res := serverTest(srv, req)

		code := 200
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Invalid Form Values", func(t *testing.T) {
		db.FindAccountsMethod = func(ids ...string) ([]waukeen.Account, error) {
			return []waukeen.Account{{ID: "1"}}, nil
		}
		db.FindTransactionsMethod = func(got waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error) {
			/* FIXME
			want := waukeen.TransactionsDBOptions{}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("wants %+v, got %+v", want, got)
			}
			*/
			return nil, nil
		}
		req := httptest.NewRequest("GET", "/accounts", nil)
		req.Form = url.Values{}
		req.Form.Set("start", "20161024")
		req.Form.Set("end", "20161024")
		req.Form.Set("transaction_type", "a")
		req.Form.Set("tags", "   ")
		res := serverTest(srv, req)

		code := 200
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Valid Form Values", func(t *testing.T) {
		db.FindAccountsMethod = func(got ...string) ([]waukeen.Account, error) {
			want := 0
			if len(got) != 0 {
				t.Errorf("wants account list to be %d length, got %d", want, got)
			}
			return []waukeen.Account{{ID: "2"}}, nil
		}
		db.FindTransactionsMethod = func(opts waukeen.TransactionsDBOptions) ([]waukeen.Transaction, error) {
			filled := waukeen.TransactionsDBOptions{
				Accounts: []string{"2"},
				Types:    []waukeen.TransactionType{waukeen.Credit},
				Start:    time.Date(2016, 10, 01, 0, 0, 0, 0, time.UTC),
				End:      time.Date(2016, 10, 31, 0, 0, 0, 0, time.UTC),
				Tags:     []string{"first", "second"},
			}
			if !reflect.DeepEqual(opts, filled) {
				t.Errorf("wants options to be %+v, got %+v", filled, opts)
			}
			return nil, nil
		}
		req := httptest.NewRequest("GET", "/accounts", nil)
		req.Form = url.Values{}
		req.Form.Set("start", "2016-10")
		req.Form.Set("end", "2016-10")
		req.Form.Set("transaction_type", "1")
		req.Form.Set("tags", "first, second ")
		req.Form.Set("account", "2")
		res := serverTest(srv, req)

		code := 200
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})
}
