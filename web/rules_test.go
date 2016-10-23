package web

import (
	"net/http/httptest"
	"testing"
)

func TestNewRule(t *testing.T) {
	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/rules/new", nil)
		res := ServerTest(nil, req)

		code := 405
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Valid Method", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/rules/new", nil)
		res := ServerTest(nil, req)

		code := 200
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})
}
