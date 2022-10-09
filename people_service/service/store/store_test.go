package store

import (
	"testing"
)

func TestNewStore(t *testing.T) {
	if _, err := NewStore("postgres://bragin:Straightein_12@95.217.232.188:7777/bragin"); err != nil {
		t.Errorf("fdsf %v", err)
	}

}

func TestListPeople(t *testing.T) {
	st, _ := NewStore("postgres://bragin:Straightein_12@95.217.232.188:7777/bragin")
	_, err := st.ListPeople()
	if err != nil {
		t.Errorf("fdsf %v", err)
	}

}

func TestGetPeopleById(t *testing.T) {
	st, _ := NewStore("postgres://bragin:Straightein_12@95.217.232.188:7777/bragin")
	_, err := st.GetPeopleByID("1")
	if err != nil {
		t.Errorf("fdsf %v", err)
	}
}
