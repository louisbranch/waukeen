package mock

import (
	"testing"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/web"
)

func TestMockSatisfiesInterfaces(t *testing.T) {
	var _ web.Template = &Template{}
	var _ waukeen.RulesImporter = &RulesImporter{}
	var _ waukeen.StatementsImporter = &StatementsImporter{}
	var _ waukeen.TransactionTransformer = &TransactionTransformer{}
	var _ waukeen.Database = &Database{}
	var _ waukeen.BudgetCalculator = &BudgetCalculator{}
}
