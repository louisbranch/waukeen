package server

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/luizbranco/waukeen/mock"
)

func serverTest(srv *Server, req *http.Request) *httptest.ResponseRecorder {
	if srv == nil {
		srv = &Server{}
	}
	if srv.Template == nil {
		tpl := &mock.Template{}

		tpl.RenderMethod = func(w io.Writer, data interface{}, path ...string) error {
			return nil
		}
		srv.Template = tpl
	}

	res := httptest.NewRecorder()
	mux := srv.NewServeMux()
	mux.ServeHTTP(res, req)

	return res
}

func fileUpload(name, uri string) *http.Request {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(name, tmpfile.Name())
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(part, tmpfile)
	if err != nil {
		log.Fatal(err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}
