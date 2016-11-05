package calc

import (
	"github.com/bradfitz/slice" // FIX Go 1.8 has built-in slice sort
	"github.com/luizbranco/waukeen"
)

type Budgeter struct{}

func (Budgeter) Calculate(months int, trs []waukeen.Transaction,
	tags []waukeen.Tag) []waukeen.Budget {

	var budget []waukeen.Budget

	m := make(map[string]waukeen.Budget)

	for _, tr := range trs {
		if len(tr.Tags) == 0 {
			tr.Tags = []string{"other"}
		}

		for _, tag := range tr.Tags {
			b := m[tag]
			b.Transactions++
			b.Spent += (tr.Amount * -1)
			b.Tag = tag
			m[tag] = b
		}
	}

	for _, t := range tags {
		b := m[t.Name]
		b.Tag = t.Name
		b.Planned = t.MonthlyBudget * int64(months)
		m[t.Name] = b
	}

	for _, b := range m {
		budget = append(budget, b)
	}

	slice.Sort(budget, func(i, j int) bool {
		return budget[i].Spent > budget[j].Spent
	})

	return budget
}
