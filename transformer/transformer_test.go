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
			waukeen.Transaction{Title: "Duty Free TORONTO", Alias: "Duty Free"},
			waukeen.Rule{Type: waukeen.ReplaceRule, Match: "toronto", Result: ""},
		},
		{
			waukeen.Transaction{Title: "MY FONE PLUS #123"},
			waukeen.Transaction{Title: "MY FONE PLUS #123", Alias: "Wind #123"},
			waukeen.Rule{Type: waukeen.ReplaceRule, Match: "MY FONE PLUS", Result: "Wind"},
		},
		{
			waukeen.Transaction{Title: "MY FONE"},
			waukeen.Transaction{Title: "MY FONE"},
			waukeen.Rule{Type: waukeen.ReplaceRule, Match: "MY FONE PLUS", Result: "Wind"},
		},
		{
			waukeen.Transaction{Title: "King Tacos?"},
			waukeen.Transaction{Title: "King Tacos?", Alias: "King Tacos"},
			waukeen.Rule{Type: waukeen.ReplaceRule, Match: "Tacos?", Result: "Tacos"},
		},
		{
			waukeen.Transaction{Title: "Pizzahut"},
			waukeen.Transaction{Title: "Pizzahut"},
			waukeen.Rule{Type: waukeen.ReplaceRule, Match: "pizza", Result: "pizza"},
		},
		{
			waukeen.Transaction{Title: "Pizzahut"},
			waukeen.Transaction{Title: "Pizzahut", Tags: []string{"pizza"}},
			waukeen.Rule{Type: waukeen.TagRule, Match: "pizzahut", Result: "pizza"},
		},
		{
			waukeen.Transaction{Title: "Pizzahut", Tags: []string{"pizza"}},
			waukeen.Transaction{Title: "Pizzahut", Tags: []string{"pizza"}},
			waukeen.Rule{Type: waukeen.TagRule, Match: "pizzahut", Result: "pizza"},
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
