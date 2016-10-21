package mock

import (
	"io"

	"github.com/luizbranco/waukeen"
)

type Template struct {
	RenderMethod func(w io.Writer, data interface{}, path ...string) error
}

func (m Template) Render(w io.Writer, data interface{}, path ...string) error {
	return m.RenderMethod(w, data, path...)
}

type StatementImporter struct {
	ImportMethod func(io.Reader) ([]waukeen.Statement, error)
}

func (m StatementImporter) Import(in io.Reader) ([]waukeen.Statement, error) {
	return m.ImportMethod(in)
}
