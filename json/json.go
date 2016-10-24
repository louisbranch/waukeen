package json

import (
	"encoding/json"
	"io"

	"github.com/luizbranco/waukeen"
)

type Rules struct{}

func (*Rules) Import(in io.Reader) ([]waukeen.Rule, error) {
	var rules []waukeen.Rule
	dec := json.NewDecoder(in)
	err := dec.Decode(&rules)
	return rules, err
}
