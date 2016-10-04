package transformer

import (
	"regexp"

	"github.com/luizbranco/waukeen"
)

type Text struct{}

func (Text) Transform(t *waukeen.Transaction, r waukeen.Rule) {
	re := regexp.MustCompile("(?i)" + r.Match)
	switch r.Type {
	case waukeen.ReplaceRule:
		t.Title = re.ReplaceAllLiteralString(t.Title, r.Result)
	case waukeen.TagRule:
		if re.MatchString(t.Title) {
			t.Tags = append(t.Tags, r.Result)
		}
	}
}
