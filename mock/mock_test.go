package mock

import (
	"testing"

	"github.com/luizbranco/waukeen"
)

func TestMockSatisfiesInterfaces(t *testing.T) {
	var _ waukeen.Template = &Template{}
	var _ waukeen.RulesImporter = &RulesImporter{}
	var _ waukeen.StatementsImporter = &StatementsImporter{}
	var _ waukeen.TransactionTransformer = &TransactionTransformer{}
	var _ waukeen.Database = &Database{}
}
