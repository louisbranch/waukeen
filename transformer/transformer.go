package transformer

import (
	"regexp"
	"strings"

	"github.com/luizbranco/waukeen"
)

type Text struct{}

func (Text) Transform(t *waukeen.Transaction, r waukeen.Rule) {
	re := regexp.MustCompile(`(?i)(^|\s)\Q` + r.Match + `\E($|\s)`)
	switch r.Type {
	case waukeen.ReplaceRule:
		alias := re.ReplaceAllString(t.Title, "${1}"+r.Result+"${2}")
		if alias != t.Title {
			t.Alias = strings.Trim(alias, " ")
		}
	case waukeen.TagRule:
		if !re.MatchString(t.Title) {
			return
		}
		t.AddTags(r.Result)
	}
}
