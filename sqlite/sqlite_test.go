package sqlite

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/luizbranco/waukeen"
	"github.com/luizbranco/waukeen/mock"
)

func testDB() (*DB, string) {
	tmpfile, err := ioutil.TempFile("", "waukeen_db")
	if err != nil {
		log.Fatal(err)
	}
	name := tmpfile.Name()
	db, err := New(name)
	if err != nil {
		log.Fatal(err)
	}
	return db, name
}

func TestCreateAccount(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	acc := &waukeen.Account{
		Number:   "12345",
		Name:     "Banking",
		Type:     waukeen.Checking,
		Currency: "CAD",
		Balance:  1000,
	}
	err := db.CreateAccount(acc)

	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if acc.ID != "1" {
		t.Errorf("wants account id 1, got %s", acc.ID)
	}

	acc = &waukeen.Account{Number: ""}
	err = db.CreateAccount(acc)
	if err == nil {
		t.Errorf("wants error, got none")
	}
}

func TestFindAccounts(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	want := []waukeen.Account{
		{
			ID:       "1",
			Number:   "12345",
			Name:     "Banking",
			Type:     waukeen.Checking,
			Currency: "CAD",
			Balance:  1000,
		},
		{
			ID:       "2",
			Number:   "67890",
			Name:     "Banking",
			Type:     waukeen.Savings,
			Currency: "USD",
			Balance:  -500,
		},
	}

	for _, a := range want {
		err := db.CreateAccount(&a)

		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}
	}

	got, err := db.FindAccounts()
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wants %+v, got %+v", want, got)
	}
}

func TestFindAccount(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	want := &waukeen.Account{
		ID:       "1",
		Number:   "12345",
		Name:     "Banking",
		Type:     waukeen.Checking,
		Currency: "CAD",
		Balance:  1000,
	}

	err := db.CreateAccount(want)

	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	got, err := db.FindAccount("12345")
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wants %+v, got %+v", want, got)
	}
}

func TestUpdateAccount(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	want := &waukeen.Account{
		ID:       "1",
		Number:   "12345",
		Name:     "Banking",
		Type:     waukeen.Checking,
		Currency: "CAD",
		Balance:  1000,
	}

	err := db.CreateAccount(want)

	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	want.Number = "02468"
	err = db.UpdateAccount(want)
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	got, err := db.FindAccount("02468")
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wants %+v, got %+v", want, got)
	}
}

func TestCreateTransaction(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	tr := &waukeen.Transaction{
		AccountID:   "1",
		FITID:       "12345",
		Type:        waukeen.Debit,
		Title:       "First Transaction",
		Alias:       "Renamed Transaction",
		Description: "Surcharge",
		Amount:      9999,
		Date:        time.Now(),
	}
	err := db.CreateTransaction(tr)

	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if tr.ID != "1" {
		t.Errorf("wants transaction id 1, got %s", tr.ID)
	}

	tr = &waukeen.Transaction{FITID: ""}
	err = db.CreateTransaction(tr)

	if err == nil {
		t.Errorf("wants error, got none")
	}
}

func TestCreateRule(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	r := &waukeen.Rule{
		AccountID: "1",
		Type:      waukeen.TagRule,
		Match:     "dominos",
		Result:    "pizza",
	}
	err := db.CreateRule(r)

	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if r.ID != "1" {
		t.Errorf("wants transaction id 1, got %s", r.ID)
	}

	r = &waukeen.Rule{Match: ""}
	err = db.CreateRule(r)

	if err == nil {
		t.Errorf("wants error, got none")
	}
}

func TestFindTransactions(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	tr1 := waukeen.Transaction{
		ID:        "1",
		AccountID: "1",
		FITID:     "01",
		Type:      waukeen.Debit,
		Title:     "1st",
		Date:      time.Date(2016, 10, 1, 0, 0, 0, 0, time.UTC),
		Tags:      []string{"groceries", "restaurants"},
	}
	tr2 := waukeen.Transaction{
		ID:        "2",
		AccountID: "1",
		FITID:     "02",
		Type:      waukeen.Credit,
		Title:     "2nd",
		Date:      time.Date(2016, 10, 5, 0, 0, 0, 0, time.UTC),
	}
	tr3 := waukeen.Transaction{
		ID:        "3",
		AccountID: "2",
		FITID:     "03",
		Type:      waukeen.Debit,
		Title:     "3rd",
		Tags:      []string{"transportation"},
		Date:      time.Date(2016, 10, 10, 0, 0, 0, 0, time.UTC),
	}
	tr4 := waukeen.Transaction{
		ID:        "4",
		AccountID: "2",
		FITID:     "04",
		Type:      waukeen.Credit,
		Title:     "4th",
		Date:      time.Date(2016, 10, 15, 0, 0, 0, 0, time.UTC),
		Tags:      []string{"groceries"},
	}

	for _, tr := range []waukeen.Transaction{tr1, tr2, tr3, tr4} {
		err := db.CreateTransaction(&tr)
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}
	}

	cases := []struct {
		opts waukeen.TransactionsDBOptions
		want []waukeen.Transaction
	}{
		{
			waukeen.TransactionsDBOptions{Accounts: []string{"1"}},
			[]waukeen.Transaction{tr1, tr2},
		},
		{
			waukeen.TransactionsDBOptions{
				Accounts: []string{"1"},
				Start:    time.Date(2016, 10, 15, 0, 0, 0, 0, time.UTC),
			},
			nil,
		},
		{
			waukeen.TransactionsDBOptions{
				Accounts: []string{"1"},
				Start:    time.Date(2016, 10, 5, 0, 0, 0, 0, time.UTC),
			},
			[]waukeen.Transaction{tr2},
		},
		{
			waukeen.TransactionsDBOptions{
				Accounts: []string{"1"},
				End:      time.Date(2016, 10, 1, 0, 0, 0, 0, time.UTC),
			},
			[]waukeen.Transaction{tr1},
		},
		{
			waukeen.TransactionsDBOptions{
				Start: time.Date(2016, 10, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2016, 10, 10, 0, 0, 0, 0, time.UTC),
			},
			[]waukeen.Transaction{tr1, tr2, tr3},
		},
		{
			waukeen.TransactionsDBOptions{
				Types: []waukeen.TransactionType{waukeen.Debit, waukeen.Credit},
			},
			[]waukeen.Transaction{tr1, tr2, tr3, tr4},
		},
		{
			waukeen.TransactionsDBOptions{
				Types: []waukeen.TransactionType{waukeen.Debit},
			},
			[]waukeen.Transaction{tr1, tr3},
		},
		{
			waukeen.TransactionsDBOptions{
				Tags: []string{"groceries"},
			},
			[]waukeen.Transaction{tr1, tr4},
		},
		{
			waukeen.TransactionsDBOptions{
				Tags: []string{"groceries", "transportation"},
			},
			[]waukeen.Transaction{tr1, tr3, tr4},
		},
		{
			waukeen.TransactionsDBOptions{
				Accounts: []string{"2"},
				Types:    []waukeen.TransactionType{waukeen.Debit},
				Start:    time.Date(2016, 10, 1, 0, 0, 0, 0, time.UTC),
				End:      time.Date(2016, 10, 10, 0, 0, 0, 0, time.UTC),
			},
			[]waukeen.Transaction{tr3},
		},
		{
			waukeen.TransactionsDBOptions{
				Accounts: []string{"1", "2"},
				Types:    []waukeen.TransactionType{waukeen.Credit, waukeen.Debit},
				Start:    time.Date(2016, 10, 1, 0, 0, 0, 0, time.UTC),
				End:      time.Date(2016, 10, 10, 0, 0, 0, 0, time.UTC),
				Tags:     []string{"restaurants", "transportation"},
			},
			[]waukeen.Transaction{tr1, tr3},
		},
	}

	for _, c := range cases {
		got, err := db.FindTransactions(c.opts)
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}
		if !reflect.DeepEqual(c.want, got) {
			t.Errorf("wants\n%+v\ngot\n%+v", c.want, got)
		}
	}
}

func TestFindRules(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	r1 := waukeen.Rule{
		AccountID: "1",
		ID:        "1",
		Type:      waukeen.TagRule,
		Match:     "dominos",
		Result:    "pizza",
	}

	r2 := waukeen.Rule{
		AccountID: "",
		ID:        "2",
		Type:      waukeen.TagRule,
		Match:     "dominos",
		Result:    "pizza",
	}

	for _, r := range []waukeen.Rule{r1, r2} {
		err := db.CreateRule(&r)
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}
	}

	want := []waukeen.Rule{r1, r2}
	got, err := db.FindRules("1")
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wants %+v, got %+v", want, got)
	}

	want = []waukeen.Rule{r2}
	got, err = db.FindRules("2")
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wants %+v, got %+v", want, got)
	}
}

func TestCreateStatement(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	transformer := &mock.TransactionTransformer{}
	transformer.TransformMethod = func(t *waukeen.Transaction, r waukeen.Rule) {
		t.Alias = r.Result
	}

	r := &waukeen.Rule{
		Type:   waukeen.ReplaceRule,
		Match:  "something",
		Result: "New Alias",
	}

	err := db.CreateRule(r)
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	t.Run("Invalid Account", func(t *testing.T) {
		stmt := waukeen.Statement{
			Account: waukeen.Account{},
		}
		err := db.CreateStatement(stmt, transformer)
		if err == nil {
			t.Errorf("wants error, got none")
		}
	})

	t.Run("Invalid Transaction", func(t *testing.T) {
		stmt := waukeen.Statement{
			Account: waukeen.Account{Number: "12345"},
			Transactions: []waukeen.Transaction{
				{},
			},
		}
		err := db.CreateStatement(stmt, transformer)
		if err == nil {
			t.Errorf("wants error, got none")
		}
	})

	t.Run("Valid Statement", func(t *testing.T) {
		stmt := waukeen.Statement{
			Account: waukeen.Account{Number: "12345"},
			Transactions: []waukeen.Transaction{
				{FITID: "67890", Title: "First"},
			},
		}
		err := db.CreateStatement(stmt, transformer)
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}

		acc, err := db.FindAccount("12345")
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}

		trs, err := db.FindTransactions(waukeen.TransactionsDBOptions{Accounts: []string{acc.ID}})
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}

		if len(trs) != 1 {
			t.Errorf("wants 1 transacton, got %+v", trs)
		}

		want := "New Alias"
		got := trs[0].Alias

		if got != want {
			t.Errorf("wants transaction alias to be %s, got %s", want, got)
		}
	})
}

func TestCreateTag(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	t.Run("Invalid Tag", func(t *testing.T) {
		tag := &waukeen.Tag{}
		err := db.CreateTag(tag)
		if err == nil {
			t.Errorf("wants error, got none")
		}
	})

	t.Run("Valid Tag", func(t *testing.T) {
		tag := &waukeen.Tag{Name: "Test"}
		err := db.CreateTag(tag)
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}
		if tag.ID != "1" {
			t.Errorf("wants tag id %s, got %s", "1", tag.ID)
		}
	})

	t.Run("Duplicated Tag", func(t *testing.T) {
		tag := &waukeen.Tag{Name: "Test"}
		err := db.CreateTag(tag)
		if err == nil {
			t.Errorf("wants error, got none")
		}
	})
}

func TestFindTags(t *testing.T) {
	db, path := testDB()
	defer os.Remove(path)

	t1 := waukeen.Tag{ID: "1", Name: "foo"}
	t2 := waukeen.Tag{ID: "2", Name: "bar"}
	t3 := waukeen.Tag{ID: "3", Name: "baz"}

	for _, tag := range []waukeen.Tag{t1, t2, t3} {
		err := db.CreateTag(&tag)
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}
	}

	t.Run("No Match", func(t *testing.T) {
		tags, err := db.FindTags("xa")
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}
		if len(tags) != 0 {
			t.Errorf("wants no tags, got %+v", tags)
		}
	})

	t.Run("Multiple Matches", func(t *testing.T) {
		want := []waukeen.Tag{t2, t3}
		got, err := db.FindTags("ba")
		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("wants %+v, got %+v", want, got)
		}
	})
}
