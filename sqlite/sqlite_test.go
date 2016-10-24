package sqlite

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/luizbranco/waukeen"
)

func TestCreateAccount(t *testing.T) {
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

	acc := &waukeen.Account{
		Number:   "12345",
		Name:     "Banking",
		Type:     waukeen.Checking,
		Currency: "CAD",
		Balance:  1000,
	}
	err = db.CreateAccount(acc)

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
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

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
		err = db.CreateAccount(&a)

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
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

	want := &waukeen.Account{
		ID:       "1",
		Number:   "12345",
		Name:     "Banking",
		Type:     waukeen.Checking,
		Currency: "CAD",
		Balance:  1000,
	}

	err = db.CreateAccount(want)

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
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

	want := &waukeen.Account{
		ID:       "1",
		Number:   "12345",
		Name:     "Banking",
		Type:     waukeen.Checking,
		Currency: "CAD",
		Balance:  1000,
	}

	err = db.CreateAccount(want)

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
