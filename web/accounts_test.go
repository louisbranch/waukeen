package web

import (
	"net/http/httptest"
	"testing"
)

func TestAccounts(t *testing.T) {
	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/accounts", nil)
		res := serverTest(nil, req)

		code := 405
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Invalid Start Date", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/accounts", nil)
		res := serverTest(nil, req)

		code := 404
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})
}
