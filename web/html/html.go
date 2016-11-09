package html

import (
	"fmt"
	"html/template"
	"io"
	"math"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type HTML struct {
	basepath string
	sync     sync.RWMutex
	cache    map[string]*template.Template
}

func New(basepath string) *HTML {
	return &HTML{
		basepath: basepath,
		cache:    make(map[string]*template.Template),
	}
}

func (h *HTML) Render(w io.Writer, data interface{}, paths ...string) error {
	for i, n := range paths {
		p := []string{h.basepath}
		p = append(p, strings.Split(n+".html", "/")...)
		paths[i] = filepath.Join(p...)
	}

	tpl, err := h.parse(paths...)
	if err != nil {
		return err
	}

	err = tpl.Execute(w, data)
	return err
}

var fns = template.FuncMap{
	"currency": currency,
	"contains": contains,
}

func (h *HTML) parse(names ...string) (tpl *template.Template, err error) {
	cp := make([]string, len(names))
	copy(cp, names)
	sort.Strings(cp)
	id := strings.Join(cp, ":")

	h.sync.RLock()
	tpl, ok := h.cache[id]
	h.sync.RUnlock()

	if !ok {
		tpl = template.New(path.Base(names[0])).Funcs(fns)

		tpl, err = tpl.ParseFiles(names...)
		if err != nil {
			return nil, err
		}
		h.sync.Lock()
		//TODO h.cache[id] = tpl
		h.sync.Unlock()
	}

	return tpl, nil
}

func currency(amount int64) string {
	symbol := "$"
	if amount < 0 {
		symbol = "-" + symbol
	}
	return fmt.Sprintf("%s%.2f", symbol, math.Abs(float64(amount))/100)
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
