package krot

import "testing"

var k = New()

func TestEat(t *testing.T) {
	k.Eat(true)
	got := k.hp
	exp := -30
	if got != exp {
		t.Errorf("got %v,but exp %v", got, exp)
	}
}

func TestDig(t *testing.T) {
	k.Dig(true)
	got := k.hp
	exp := -30
	if got != exp {
		t.Errorf("got %v,but exp %v", got, exp)
	}
}

func TestRep_changes(t *testing.T) {
	got := repChanges(30, 50)
	exp := 5
	if got != exp {
		t.Errorf("got %v,but exp %v", got, exp)
	}
}

func TestSleep(t *testing.T) {
	k.Sleep()
	got := k.noraLen
	exp := -2
	if got != exp {
		t.Errorf("got %v,but exp %v", got, exp)
	}
}
