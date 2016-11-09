package web

import (
	"io"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/pkg/errors"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/mock"
)

func TestNewRule(t *testing.T) {
	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/rules/new", nil)
		res := serverTest(nil, req)

		code := 405
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Valid Method", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/rules/new", nil)
		res := serverTest(nil, req)

		code := 200
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})
}

func TestRules(t *testing.T) {
	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/rules/", nil)
		res := serverTest(nil, req)

		code := 405
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Get rules DB error", func(t *testing.T) {
		db := &mock.Database{}
		db.FindRulesMethod = func(...string) ([]waukeen.Rule, error) {
			return nil, errors.New("not implemented")
		}
		srv := &Server{DB: db}

		req := httptest.NewRequest("GET", "/rules/", nil)
		res := serverTest(srv, req)

		code := 500
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Get rules", func(t *testing.T) {
		db := &mock.Database{}
		db.FindRulesMethod = func(...string) ([]waukeen.Rule, error) {
			return nil, nil
		}
		srv := &Server{DB: db}

		req := httptest.NewRequest("GET", "/rules/", nil)
		res := serverTest(srv, req)

		code := 200
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Post new rule invalid type", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/rules/", nil)
		req.Form = url.Values{}
		req.Form.Set("type", "a")

		res := serverTest(nil, req)

		code := 500
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Post new rule invalid rule", func(t *testing.T) {
		db := &mock.Database{}
		db.CreateRuleMethod = func(*waukeen.Rule) error {
			return errors.New("not implemented")
		}
		srv := &Server{DB: db}

		req := httptest.NewRequest("POST", "/rules/", nil)
		req.Form = url.Values{}
		req.Form.Set("type", "1")

		res := serverTest(srv, req)

		code := 500
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Post new rule valid rule", func(t *testing.T) {
		rule := &waukeen.Rule{
			Type:   waukeen.TagRule,
			Match:  "dominos",
			Result: "pizza",
		}

		db := &mock.Database{}
		db.CreateRuleMethod = func(r *waukeen.Rule) error {
			if !reflect.DeepEqual(r, rule) {
				t.Errorf("want %s, got %s", rule, r)
			}
			return nil
		}
		srv := &Server{DB: db}

		req := httptest.NewRequest("POST", "/rules/", nil)
		req.Form = url.Values{}
		req.Form.Set("type", strconv.Itoa(int(rule.Type)))
		req.Form.Set("match", rule.Match)
		req.Form.Set("result", rule.Result)

		res := serverTest(srv, req)

		code := 302
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}

		url := "/rules/"
		loc := res.Header().Get("Location")

		if url != loc {
			t.Errorf("wants %s redirect url, got %s", url, loc)
		}
	})
}

func TestImportRules(t *testing.T) {
	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/rules/import", nil)
		res := serverTest(nil, req)

		code := 405
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Get rules", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/rules/import", nil)
		res := serverTest(nil, req)

		code := 200
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Post rule file error", func(t *testing.T) {
		importer := &mock.RulesImporter{}
		importer.ImportMethod = func(io.Reader) ([]waukeen.Rule, error) {
			return nil, errors.New("not implemented")
		}
		srv := &Server{RulesImporter: importer}

		req := fileUpload("rules", "/rules/import")
		res := serverTest(srv, req)

		code := 500
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Post rule file db error", func(t *testing.T) {
		importer := &mock.RulesImporter{}
		importer.ImportMethod = func(io.Reader) ([]waukeen.Rule, error) {
			return []waukeen.Rule{waukeen.Rule{}}, nil
		}
		db := &mock.Database{}
		db.CreateRuleMethod = func(r *waukeen.Rule) error {
			return errors.New("not implemented")
		}
		srv := &Server{RulesImporter: importer, DB: db}

		req := fileUpload("rules", "/rules/import")
		res := serverTest(srv, req)

		code := 500
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Post valid rules", func(t *testing.T) {
		rule := &waukeen.Rule{
			Type:   waukeen.TagRule,
			Match:  "dominos",
			Result: "pizza",
		}
		importer := &mock.RulesImporter{}
		importer.ImportMethod = func(io.Reader) ([]waukeen.Rule, error) {
			return []waukeen.Rule{*rule}, nil
		}
		db := &mock.Database{}
		db.CreateRuleMethod = func(r *waukeen.Rule) error {
			if !reflect.DeepEqual(r, rule) {
				t.Errorf("want %s, got %s", rule, r)
			}
			return nil
		}
		srv := &Server{RulesImporter: importer, DB: db}

		req := fileUpload("rules", "/rules/import")
		res := serverTest(srv, req)

		code := 302
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}

		url := "/rules"
		loc := res.Header().Get("Location")

		if url != loc {
			t.Errorf("wants %s redirect url, got %s", url, loc)
		}
	})
}
