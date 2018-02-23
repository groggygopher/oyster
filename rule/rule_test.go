package rule

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/groggygopher/oyster/register"
)

const (
	amount = 13.37
)

func pointer(f float64) *float64 {
	return &f
}

var (
	now    = time.Now()
	before = now.Add(-60 * time.Minute)
	after  = now.Add(60 * time.Minute)

	trans = &register.Transaction{
		Date:        &now,
		Description: "test",
		Amount:      amount,
	}

	descriptionMatch   = &Rule{Description: &Description{r: regexp.MustCompile("test")}}
	descriptionNoMatch = &Rule{Description: &Description{r: regexp.MustCompile("bad")}}

	dateBeforeMatch   = &Rule{DateBetween: &DateRange{Before: &after}}
	dateBeforeNoMatch = &Rule{DateBetween: &DateRange{Before: &before}}
	dateAfterMatch    = &Rule{DateBetween: &DateRange{After: &before}}
	dateAfterNoMatch  = &Rule{DateBetween: &DateRange{After: &after}}

	amountMinMatch   = &Rule{AmountBetween: &AmountRange{Min: pointer(0)}}
	amountMinNoMatch = &Rule{AmountBetween: &AmountRange{Min: pointer(100)}}
	amountMaxMatch   = &Rule{AmountBetween: &AmountRange{Max: pointer(100)}}
	amountMaxNoMatch = &Rule{AmountBetween: &AmountRange{Max: pointer(0)}}
)

func TestRuleEvaluate(t *testing.T) {
	tests := []struct {
		label string
		rule  *Rule
		want  bool
	}{
		{
			label: "match description",
			rule:  descriptionMatch,
			want:  true,
		},
		{
			label: "no match description",
			rule:  descriptionNoMatch,
			want:  false,
		},
		{
			label: "match date before",
			rule:  dateBeforeMatch,
			want:  true,
		},
		{
			label: "no match date before",
			rule:  dateBeforeNoMatch,
			want:  false,
		},
		{
			label: "match date after",
			rule:  dateAfterMatch,
			want:  true,
		},
		{
			label: "no match date after",
			rule:  dateAfterNoMatch,
			want:  false,
		},
		{
			label: "match amount min",
			rule:  amountMinMatch,
			want:  true,
		},
		{
			label: "no match amount min",
			rule:  amountMinNoMatch,
			want:  false,
		},
		{
			label: "match amount max",
			rule:  amountMaxMatch,
			want:  true,
		},
		{
			label: "no match amount max",
			rule:  amountMaxNoMatch,
			want:  false,
		},
		{
			label: "and match",
			rule: &Rule{
				And: []*Rule{amountMinMatch, amountMaxMatch},
			},
			want: true,
		},
		{
			label: "and no match",
			rule: &Rule{
				And: []*Rule{amountMinMatch, amountMinNoMatch},
			},
			want: false,
		},
		{
			label: "or match",
			rule: &Rule{
				Or: []*Rule{amountMinMatch, amountMaxNoMatch},
			},
			want: true,
		},
		{
			label: "or no match",
			rule: &Rule{
				Or: []*Rule{amountMinNoMatch, amountMaxNoMatch},
			},
			want: false,
		},
		{
			label: "match local and And",
			rule: &Rule{
				Description: &Description{r: regexp.MustCompile("test")},
				And:         []*Rule{amountMinMatch},
			},
			want: true,
		},
		{
			label: "match no local and And",
			rule: &Rule{
				Description: &Description{r: regexp.MustCompile("bad")},
				And:         []*Rule{amountMinMatch},
			},
			want: false,
		},
		{
			label: "match local and no And",
			rule: &Rule{
				Description: &Description{r: regexp.MustCompile("test")},
				And:         []*Rule{amountMinNoMatch},
			},
			want: false,
		},
		{
			label: "match no local and no And",
			rule: &Rule{
				Description: &Description{r: regexp.MustCompile("bad")},
				And:         []*Rule{amountMinNoMatch},
			},
			want: false,
		},
		{
			label: "match local and Or",
			rule: &Rule{
				Description: &Description{r: regexp.MustCompile("test")},
				Or:          []*Rule{amountMinMatch},
			},
			want: true,
		},
		{
			label: "match no local and Or",
			rule: &Rule{
				Description: &Description{r: regexp.MustCompile("bad")},
				Or:          []*Rule{amountMinMatch},
			},
			want: true,
		},
		{
			label: "match local and no Or",
			rule: &Rule{
				Description: &Description{r: regexp.MustCompile("test")},
				Or:          []*Rule{amountMinNoMatch},
			},
			want: true,
		},
		{
			label: "match no local and no Or",
			rule: &Rule{
				Description: &Description{r: regexp.MustCompile("bad")},
				Or:          []*Rule{amountMinNoMatch},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			if got, want := test.rule.Evaluate(trans), test.want; got != want {
				t.Errorf("got: %t, want: %t", got, want)
			}
		})
	}
}

func TestDescriptionJSON(t *testing.T) {
	const re = "hello.world"
	d := &Description{r: regexp.MustCompile(re)}

	jsonBytes, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}

	unmarshalled := &Description{}
	if err := json.Unmarshal(jsonBytes, unmarshalled); err != nil {
		t.Fatal(err)
	}

	if got, want := unmarshalled.r.String(), re; got != want {
		t.Errorf("json encode error: got: %s, want: %s", got, want)
	}
}
