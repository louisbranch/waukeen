package html

import (
	"fmt"
	"html/template"
	"io"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/luizbranco/waukeen/web"
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

func (h *HTML) Render(w io.Writer, page web.Page) error {
	paths := append([]string{page.Layout}, page.Partials...)

	for i, n := range paths {
		p := []string{h.basepath}
		p = append(p, strings.Split(n+".html", "/")...)
		paths[i] = filepath.Join(p...)
	}

	tpl, err := h.parse(paths...)
	if err != nil {
		return err
	}

	err = tpl.Execute(w, page)
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

func currency(val int64) string {
	symbol := "$"
	if val < 0 {
		symbol = "-" + symbol
		val *= -1
	}

	res := fmt.Sprintf(".%02d", val%100)

	val = val / 100

	for val >= 1000 {
		res = fmt.Sprintf(",%03d", val%1000) + res
		val = val / 1000
	}

	res = fmt.Sprintf("%s%d", symbol, val) + res

	return res
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
