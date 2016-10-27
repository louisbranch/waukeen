package xml

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/luizbranco/waukeen"
)

func TestStatementsImport(t *testing.T) {
	importer := &Statement{}

	t.Run("Invalid XML", func(t *testing.T) {
		in := strings.NewReader(``)
		_, err := importer.Import(in)
		if err == nil {
			t.Error("wants error, got none")
		}
	})

	t.Run("Valid JSON", func(t *testing.T) {
		in := strings.NewReader(`
OFXHEADER:100
DATA:OFXSGML
VERSION:102
SECURITY:NONE
ENCODING:USASCII
CHARSET:1252
COMPRESSION:NONE
OLDFILEUID:NONE
NEWFILEUID:NONE

<OFX>
	<SIGNONMSGSRSV1>
		<SONRS>
			<STATUS>
				<CODE>0
				<SEVERITY>INFO
				<MESSAGE>OK
			</STATUS>
			<DTSERVER>20160920110857.000[-4:EDT]
			<LANGUAGE>ENG
			<INTU.BID>00025
		</SONRS>
	</SIGNONMSGSRSV1>
	<BANKMSGSRSV1>
		<STMTTRNRS>
			<TRNUID>20160920110535.000[-4:EDT]
			<STATUS>
				<CODE>0
				<SEVERITY>INFO
				<MESSAGE>OK
			</STATUS>
			<STMTRS>
				<CURDEF>CAD
				<BANKACCTFROM>
					<BANKID>170000100
					<ACCTID>1234567890
					<ACCTTYPE>SAVINGS
				</BANKACCTFROM>
				<BANKTRANLIST>
					<DTSTART>20160910
					<DTEND>20160920
					<STMTTRN>
						<TRNTYPE>DEBIT
						<DTPOSTED>20160910120000
						<TRNAMT>-49.77
						<FITID>12345
						<NAME>Credit Card Payment
						<MEMO>Customer Transfer
					</STMTTRN>
					<STMTTRN>
						<TRNTYPE>CREDIT
						<DTPOSTED>20160909120000
						<TRNAMT>49.77
						<FITID>23456
						<NAME>CANADA
						<MEMO>Provincial Payment
					</STMTTRN>
				</BANKTRANLIST>
				<LEDGERBAL>
					<BALAMT>1200.00
					<DTASOF>20160920110535.000[-4:EDT]
				</LEDGERBAL>
				<AVAILBAL>
					<BALAMT>1200.00
					<DTASOF>20160920110535.000[-4:EDT]
				</AVAILBAL>
			</STMTRS>
		</STMTTRNRS>
	</BANKMSGSRSV1>
	<CREDITCARDMSGSRSV1>
		<CCSTMTTRNRS>
			<TRNUID>20160920110857.000[-4:EDT]
			<STATUS>
				<CODE>0
				<SEVERITY>INFO
				<MESSAGE>OK
			</STATUS>
			<CCSTMTRS>
				<CURDEF>CAD
				<CCACCTFROM>
					<ACCTID>1234567890123456
				</CCACCTFROM>
				<BANKTRANLIST>
				<DTSTART>20160809
				<DTEND>20160920
				<STMTTRN>
					<TRNTYPE>DEBIT
					<DTPOSTED>20160809120000
					<TRNAMT>-163.25
					<FITID>12345
					<NAME>Game of Thrones
				</STMTTRN>
				<STMTTRN>
					<TRNTYPE>DEBIT
					<DTPOSTED>20160810120000
					<TRNAMT>-71.19
					<FITID>23456
					<NAME>Aquarium
				</STMTTRN>
			</BANKTRANLIST>
			<LEDGERBAL>
				<BALAMT>-436.14
				<DTASOF>20160920110857.000[-4:EDT]
			</LEDGERBAL>
			<AVAILBAL>
				<BALAMT>-436.14
				<DTASOF>20160920110857.000[-4:EDT]
			</AVAILBAL>
			</CCSTMTRS>
		</CCSTMTTRNRS>
	</CREDITCARDMSGSRSV1>
</OFX>
		`)
		want := []waukeen.Statement{
			{
				Account: waukeen.Account{
					Number:   "1234567890",
					Type:     waukeen.Savings,
					Currency: "CAD",
					Balance:  120000,
				},
				Transactions: []waukeen.Transaction{
					{
						FITID:       "12345",
						Type:        waukeen.Debit,
						Title:       "Credit Card Payment",
						Description: "Customer Transfer",
						Amount:      -4977,
						Date:        time.Date(2016, 9, 10, 12, 0, 0, 0, time.UTC),
					},
					{
						FITID:       "23456",
						Type:        waukeen.Credit,
						Title:       "CANADA",
						Description: "Provincial Payment",
						Amount:      4977,
						Date:        time.Date(2016, 9, 9, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			{
				Account: waukeen.Account{
					Number:   "1234567890123456",
					Type:     waukeen.CreditCard,
					Currency: "CAD",
					Balance:  -43614,
				},
				Transactions: []waukeen.Transaction{
					{
						FITID:  "12345",
						Type:   waukeen.Debit,
						Title:  "Game of Thrones",
						Amount: -16325,
						Date:   time.Date(2016, 8, 9, 12, 0, 0, 0, time.UTC),
					},
					{
						FITID:  "23456",
						Type:   waukeen.Debit,
						Title:  "Aquarium",
						Amount: -7119,
						Date:   time.Date(2016, 8, 10, 12, 0, 0, 0, time.UTC),
					},
				},
			},
		}
		got, err := importer.Import(in)
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("wants\n%+v\ngot\n%+v", want, got)
		}
	})
}
