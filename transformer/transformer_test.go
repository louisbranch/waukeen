package transformer

import (
	"reflect"
	"testing"

	"github.com/luizbranco/waukeen"
)

func TestTextTransform(t *testing.T) {
	txt := Text{}

	egs := []struct {
		in   waukeen.Transaction
		out  waukeen.Transaction
		rule waukeen.Rule
	}{
		{
			waukeen.Transaction{Title: "Duty Free TORONTO"},
			waukeen.Transaction{Title: "Duty Free"},
			waukeen.Rule{Type: waukeen.ReplaceRule, Match: " toronto", Result: ""},
		},
		{
			waukeen.Transaction{Title: "MY FONE PLUS"},
			waukeen.Transaction{Title: "Wind"},
			waukeen.Rule{Type: waukeen.ReplaceRule, Match: "MY FONE PLUS", Result: "Wind"},
		},
		{
			waukeen.Transaction{Title: "MY FONE"},
			waukeen.Transaction{Title: "MY FONE"},
			waukeen.Rule{Type: waukeen.ReplaceRule, Match: "MY FONE PLUS", Result: "Wind"},
		},
		{
			waukeen.Transaction{Title: "Pizza Pizza"},
			waukeen.Transaction{Title: "Pizza Pizza", Tags: []string{"food"}},
			waukeen.Rule{Type: waukeen.TagRule, Match: "pizza", Result: "food"},
		},
	}

	for _, eg := range egs {
		txt.Transform(&eg.in, eg.rule)
		want := &eg.out
		got := &eg.in

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}
