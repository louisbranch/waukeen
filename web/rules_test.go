package web

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/luizbranco/waukeen/mock"
)

func TestNewRule(t *testing.T) {
	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/rules/new", nil)
		res := httptest.NewRecorder()

		srv := &Server{}
		srv.newRule(res, req)

		code := 405
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})

	t.Run("Valid Method", func(t *testing.T) {
		tpl := &mock.Template{}

		tpl.RenderMethod = func(w io.Writer, data interface{}, path ...string) error {
			return nil
		}
		req := httptest.NewRequest("GET", "/rules/new", nil)
		res := httptest.NewRecorder()

		srv := &Server{Template: tpl}
		srv.newRule(res, req)

		code := 200
		if res.Code != code {
			t.Errorf("wants %d status code, got %d", code, res.Code)
		}
	})
}
