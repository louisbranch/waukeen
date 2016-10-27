package json

import (
	"reflect"
	"strings"
	"testing"

	"github.com/luizbranco/waukeen"
)

func TestRulesImport(t *testing.T) {
	importer := &Rules{}

	t.Run("Invalid JSON", func(t *testing.T) {
		in := strings.NewReader(``)
		_, err := importer.Import(in)
		if err == nil {
			t.Error("wants error, got none")
		}
	})

	t.Run("Valid JSON", func(t *testing.T) {
		in := strings.NewReader(`[
			{"type": "tag", "match": "dominos", "result": "pizza"},
			{"type": "replace", "match": "toronto", "result": "local"}
		]`)
		want := []waukeen.Rule{
			{Type: waukeen.TagRule, Match: "dominos", Result: "pizza"},
			{Type: waukeen.ReplaceRule, Match: "toronto", Result: "local"},
		}
		got, err := importer.Import(in)
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("wants %+v, got %+v", want, got)
		}
	})
}
