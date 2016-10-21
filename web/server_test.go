package web

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/mock"
)

func TestNewStatement(t *testing.T) {
	tpl := &mock.Template{}
	srv := &Server{Template: tpl}

	req := httptest.NewRequest("GET", "/statements/new", nil)
	res := httptest.NewRecorder()

	tpl.RenderMethod = func(io.Writer, interface{}, ...string) error {
		return nil
	}

	srv.newStatement(res, req)

	code := 200
	if res.Code != code {
		t.Errorf("wants %d status code, got %d", code, res.Code)
	}

	req = httptest.NewRequest("POST", "/statements/new", nil)
	res = httptest.NewRecorder()

	srv.newStatement(res, req)

	code = 405
	if res.Code != code {
		t.Errorf("wants %d status code, got %d", code, res.Code)
	}
}

func TestCreateStatement(t *testing.T) {
	importer := &mock.StatementImporter{}

	// Invalid Method
	srv := &Server{Statement: importer}

	req := httptest.NewRequest("GET", "/statements", nil)
	res := httptest.NewRecorder()

	srv.createStatement(res, req)

	code := 405
	if res.Code != code {
		t.Errorf("wants %d status code, got %d", code, res.Code)
	}

	// Missing File
	importer.ImportMethod = func(io.Reader) ([]waukeen.Statement, error) {
		return nil, errors.New("not implemented")
	}
	req = httptest.NewRequest("POST", "/statements", nil)
	res = httptest.NewRecorder()

	srv.createStatement(res, req)

	code = 400
	if res.Code != code {
		t.Errorf("wants %d status code, got %d", code, res.Code)
	}

	// Empty File

	importer.ImportMethod = func(io.Reader) ([]waukeen.Statement, error) {
		return nil, nil
	}

	req, err := fileUpload("statement", "/statements", "../mock/cc.ofx")
	if err != nil {
		t.Error(err)
	}
	res = httptest.NewRecorder()

	srv.createStatement(res, req)

	code = 200
	if res.Code != code {
		t.Errorf("wants %d status code, got %d (%s)", code, res.Code, res.Body)
	}
}

func fileUpload(name, uri, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}
