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
		t.Title = re.ReplaceAllString(t.Title, "${1}"+r.Result+"${2}")
		t.Title = strings.Trim(t.Title, " ")
	case waukeen.TagRule:
		if re.MatchString(t.Title) {
			t.Tags = append(t.Tags, r.Result)
		}
	}
}
