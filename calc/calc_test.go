package calc

import (
	"reflect"
	"testing"

	"github.com/luizbranco/waukeen"
)

func TestBudgetCalculatorInterface(t *testing.T) {
	var _ waukeen.BudgetCalculator = Budgeter{}
}

func TestCalculate(t *testing.T) {
	b := Budgeter{}

	testCases := []struct {
		months       int
		transactions []waukeen.Transaction
		tags         []waukeen.Tag
		budget       []waukeen.Budget
	}{
		{}, // empty case
		{
			months: 1,
			transactions: []waukeen.Transaction{
				{Amount: -1000, Tags: []string{"food", "pizza"}},
				{Amount: -5000, Tags: []string{"travel"}},
				{Amount: -4500, Tags: []string{"gift", "pizza"}},
			},
			tags: []waukeen.Tag{
				{Name: "food", MonthlyBudget: 2000},
				{Name: "pizza", MonthlyBudget: 5000},
			},
			budget: []waukeen.Budget{
				{Tag: "food", Transactions: 1, Planned: 2000, Spent: 1000},
				{Tag: "gift", Transactions: 1, Planned: 0, Spent: 4500},
				{Tag: "pizza", Transactions: 2, Planned: 5000, Spent: 5500},
				{Tag: "travel", Transactions: 1, Planned: 0, Spent: 5000},
			},
		},
		{
			months: 3,
			transactions: []waukeen.Transaction{
				{Amount: -1000, Tags: []string{"food"}},
				{Amount: -5000},
				{Amount: -250},
			},
			tags: []waukeen.Tag{
				{Name: "food", MonthlyBudget: 2000},
			},
			budget: []waukeen.Budget{
				{Tag: "food", Transactions: 1, Planned: 6000, Spent: 1000},
				{Tag: "other", Transactions: 2, Planned: 0, Spent: 5250},
			},
		},
	}

	for _, tc := range testCases {
		want := tc.budget
		got := b.Calculate(tc.months, tc.transactions, tc.tags)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("wants %+v, got %+v", want, got)
		}
	}
}
