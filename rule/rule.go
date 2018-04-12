package rule

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/groggygopher/oyster/register"
)

// DateRange is a date range between which a Rule can be evaluated.
type DateRange struct {
	After  *time.Time `json:"after"`
	Before *time.Time `json:"before"`
}

// AmountRange is a dollar amount range between which a Rule can be evaluated.
type AmountRange struct {
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`
}

// Description is the Transaction description RE on which a Rule can be evaluated.
type Description struct {
	r *regexp.Regexp
}

// UnmarshalJSON converts the given string to a regexp inside Description.
func (d *Description) UnmarshalJSON(b []byte) error {
	re, err := regexp.Compile(strings.Trim(string(b), `"`))
	if err != nil {
		return fmt.Errorf("not a valid regex: %s, err: %v", b, err)
	}
	*d = Description{r: re}
	return nil
}

// MarshalJSON dumps the Description regexp as its pattern string.
func (d *Description) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", d.r.String())), nil
}

// Rule is a single rule that a Transaction can be evaluated against in order to assign it to its
// Category. Client code should use a Manager to set the Category of a Transaction when it matches
// a Rule.
type Rule struct {
	Name     string `json:"name"`
	Category string `json:"category"`

	And           []*Rule      `json:"and"`
	Or            []*Rule      `json:"or"`
	Description   *Description `json:"description"`
	DateBetween   *DateRange   `json:"dateBetween"`
	AmountBetween *AmountRange `json:"amountBetween"`
}

// Evaluate checks if the Transaction matches this rule and returns true if so.
func (r *Rule) Evaluate(t *register.Transaction) bool {
	local := true
	set := false
	if r.Description != nil && t.Description != "" {
		set = true
		local = local && r.Description.r.MatchString(t.Description)
	}
	if r.DateBetween != nil && t.Date != nil {
		set = true
		isBetween := true
		if r.DateBetween.Before != nil {
			isBetween = isBetween && t.Date.Before(*r.DateBetween.Before)
		}
		if r.DateBetween.After != nil {
			isBetween = isBetween && t.Date.After(*r.DateBetween.After)
		}
		local = local && isBetween
	}
	if r.AmountBetween != nil {
		isBetween := true
		if r.AmountBetween.Min != nil {
			set = true
			isBetween = isBetween && t.Amount >= *r.AmountBetween.Min
		}
		if r.AmountBetween.Max != nil {
			set = true
			isBetween = isBetween && t.Amount <= *r.AmountBetween.Max
		}
		local = local && isBetween
	}
	if len(r.And) > 0 {
		and := true
		for _, r := range r.And {
			and = and && r.Evaluate(t)
		}
		local = local && and
	}
	if len(r.Or) > 0 {
		or := false
		for _, r := range r.Or {
			or = or || r.Evaluate(t)
		}
		// Make sure that a default true local doesn't confuse the Or logic when the rule is only an Or.
		if !set {
			local = or
		} else {
			local = local || or
		}
	}
	return local
}
