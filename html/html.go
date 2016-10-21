package html

import (
	"html/template"
	"io"
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

func (h *HTML) parse(names ...string) (tpl *template.Template, err error) {
	cp := make([]string, len(names))
	copy(cp, names)
	sort.Strings(cp)
	id := strings.Join(cp, ":")

	h.sync.RLock()
	tpl, ok := h.cache[id]
	h.sync.RUnlock()

	if !ok {
		tpl, err = template.ParseFiles(names...)
		if err != nil {
			return nil, err
		}
		h.sync.Lock()
		h.cache[id] = tpl
		h.sync.Unlock()
	}

	return tpl, nil
}
