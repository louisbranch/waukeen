package server

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/mock"
)

func TestNewStatement(t *testing.T) {
	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/statements/new", nil)
		res := serverTest(nil, req)

		code := 405
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Valid Method", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/statements/new", nil)
		res := serverTest(nil, req)

		code := 200
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})
}

func TestCreateStatement(t *testing.T) {
	importer := &mock.StatementsImporter{}
	db := &mock.Database{}

	srv := &Server{
		StatementsImporter: importer,
		DB:                 db,
	}

	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/statements", nil)
		res := serverTest(srv, req)

		code := 405
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Missing File", func(t *testing.T) {
		importer.ImportMethod = func(io.Reader) ([]waukeen.Statement, error) {
			return nil, errors.New("not implemented")
		}
		req := httptest.NewRequest("POST", "/statements", nil)
		res := serverTest(srv, req)

		code := 500
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Empty File", func(t *testing.T) {
		importer.ImportMethod = func(io.Reader) ([]waukeen.Statement, error) {
			return nil, nil
		}
		req := fileUpload("statement", "/statements")
		res := serverTest(srv, req)

		code := 500
		if res.Code != code {
			t.Errorf("wants %d status code, got %d (%s)", code, res.Code, res.Body)
		}
	})

	t.Run("Statement Error", func(t *testing.T) {
		importer.ImportMethod = func(io.Reader) ([]waukeen.Statement, error) {
			return []waukeen.Statement{{
				Account: waukeen.Account{Number: "12345"},
			}}, nil
		}

		db.CreateStatementMethod = func(waukeen.Statement, waukeen.TransactionTransformer) error {
			return errors.New("account not found")
		}

		req := fileUpload("statement", "/statements")
		res := serverTest(srv, req)

		code := 500
		if res.Code != code {
			t.Errorf("wants %d status code, got %d (%s)", code, res.Code, res.Body)
		}
	})

	t.Run("Statement(s) Successfully Imported", func(t *testing.T) {
		importer.ImportMethod = func(io.Reader) ([]waukeen.Statement, error) {
			return []waukeen.Statement{{
				Account: waukeen.Account{Number: "12345"},
			}}, nil
		}

		db.CreateStatementMethod = func(waukeen.Statement, waukeen.TransactionTransformer) error {
			return nil
		}

		req := fileUpload("statement", "/statements")
		res := serverTest(srv, req)

		code := 302
		if res.Code != code {
			t.Errorf("wants %d status code, got %d (%s)", code, res.Code, res.Body)
		}

		url := "/accounts"
		loc := res.Header().Get("Location")

		if url != loc {
			t.Errorf("wants %s redirect url, got %s", url, loc)
		}
	})
}
