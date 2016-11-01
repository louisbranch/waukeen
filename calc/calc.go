package calc

import (
	"strings"

	"github.com/bradfitz/slice" // FIX Go 1.8 has built-in slice sort
	"github.com/luizbranco/waukeen"
)

type Budgeter struct{}

func (*Budgeter) Calculate(trs []waukeen.Transaction, tags []waukeen.Tag) []waukeen.Budget {

	var budget []waukeen.Budget

	m := make(map[string]waukeen.Budget)

	for _, tr := range trs {
		for _, tag := range tr.Tags {
			b, ok := m[tag]
			if ok {
				b.Transactions += 1
				b.Spent += tr.Amount
			} else {
				b = waukeen.Budget{Tag: tag, Transactions: 1, Spent: tr.Amount}
			}
			m[tag] = b
		}
	}

	for _, t := range tags {
		b, ok := m[t.Name]
		if !ok {
			continue
		}
		b.Planned = t.Budget
		m[t.Name] = b
	}

	for _, b := range m {
		budget = append(budget, b)
	}

	slice.Sort(budget, func(i, j int) bool {
		return strings.Compare(budget[i].Tag, budget[j].Tag) < 0
	})

	return budget
}
