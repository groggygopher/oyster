package register

import (
	"fmt"
	"time"
)

func monthsKey(m time.Month, y int) string {
	return fmt.Sprintf("%d/%s", y, m)
}

// Category is a Transaction's entry against a Budget. The name of a Category should match a Budget.
type Category struct {
	Name   string
	Amount float64
}

// Budget is a simple representation of a single Budget line item across many months or years.
type Budget struct {
	id     string
	name   string
	months map[string]float64
}

// NewBudget returns a new empty Budget with the given name.
func NewBudget(name string) *Budget {
	return &Budget{
		id:     fmt.Sprintf("BUD-%s-%d", name, time.Now().Unix()),
		name:   name,
		months: make(map[string]float64),
	}
}

// Name returns the name of this Budget.
func (b *Budget) Name() string {
	return b.name
}

// SetAmount sets the amount of a Budget in the given month and year.
func (b *Budget) SetAmount(month time.Month, year int, amount float64) {
	b.months[monthsKey(month, year)] = amount
}

// Amount returns the set amount of a Budget in the given month and year. True is returned if the
// Budget was set for that month.
func (b *Budget) Amount(month time.Month, year int) (float64, bool) {
	amt, ok := b.months[monthsKey(month, year)]
	return amt, ok
}
