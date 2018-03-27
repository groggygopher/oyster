package rule

import (
	"regexp"
	"testing"

	"github.com/groggygopher/oyster/register"
)

func TestManagerEvaluate(t *testing.T) {
	trans := &register.Transaction{
		Description: "test",
	}
	tests := []struct {
		label      string
		rules      []*Rule
		wantChange bool
		wantErr    bool
	}{
		{
			label: "match single rule",
			rules: []*Rule{
				{
					Description: &Description{regexp.MustCompile("test")},
				},
			},
			wantChange: true,
		},
		{
			label: "match no rules",
			rules: []*Rule{
				{
					Description: &Description{regexp.MustCompile("nomatch")},
				},
			},
		},
		{
			label: "match multiple rules",
			rules: []*Rule{
				{
					Description: &Description{regexp.MustCompile("test")},
				},
				{
					Description: &Description{regexp.MustCompile("test")},
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			mngr := NewEmptyManager()
			for _, r := range test.rules {
				mngr.AddRule(r)
			}
			change, err := mngr.Evaluate(trans)
			if got, want := change, test.wantChange; got != want {
				t.Errorf("change: got: %t, want: %t", got, want)
			}
			if got, want := err != nil, test.wantErr; got != want {
				t.Errorf("error: got: %t, want: %t, err: %v", got, want, err)
			}
		})
	}
}

func TestManagerEvaluate_SetCategory(t *testing.T) {
	const category = "test"
	const amount = 13.37
	trans := &register.Transaction{
		Description: "test",
		Amount:      amount,
	}
	rule := &Rule{
		Category:    category,
		Description: &Description{regexp.MustCompile("test")},
	}
	mngr := NewEmptyManager()
	mngr.AddRule(rule)
	changed, err := mngr.Evaluate(trans)
	if err != nil {
		t.Fatalf("manager.Evaluate: %v", err)
	}
	if got, want := changed, true; got != want {
		t.Fatalf("changed: got: %t, want: %t", got, want)
	}

	if got, want := len(trans.Category), 1; got != want {
		t.Fatalf("transaction.Category length: got: %d, want: %d", got, want)
	}
	cat := trans.Category[0]
	if got, want := cat.Name, category; got != want {
		t.Errorf("category name: got: %s, want: %s", got, want)
	}
	if got, want := cat.Amount, amount; got != want {
		t.Errorf("category amount: got: %f, want: %f", got, want)
	}
}
