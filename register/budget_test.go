package register

import (
	"testing"
	"time"
)

func TestSetBudget(t *testing.T) {
	b := NewBudget("test")
	if _, ok := b.Amount(time.January, 2018); ok {
		t.Error("new budget should not have any entries")
	}

	b.SetAmount(time.December, 2018, 13.37)
	if _, ok := b.Amount(time.January, 2018); ok {
		t.Error("Jan 2018 should not have any entries")
	}

	amt, ok := b.Amount(time.December, 2018)
	if !ok {
		t.Error("Dec 2018 should have an entry")
	}
	if got, want := amt, 13.37; got != want {
		t.Errorf("amount, got: %f, want: %f", got, want)
	}
}
