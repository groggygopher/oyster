package rule

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/groggygopher/oyster/register"
)

// Manager manages the evaluation of a Transaction against zero or more rules.
type Manager struct {
	rules []*Rule
}

// NewEmptyManager returns a new empty rule Manager.
func NewEmptyManager() *Manager {
	return &Manager{
		rules: []*Rule{},
	}
}

// NewManager returns a rule Manager with the given rules.
func NewManager(rs []*Rule) *Manager {
	return &Manager{
		rules: rs,
	}
}

// Rules returns a slice of this manager's rules.
func (m *Manager) Rules() []*Rule {
	if m == nil {
		return nil
	}
	return m.rules
}

// AddRule adds a rule to this Manager.
func (m *Manager) AddRule(r *Rule) {
	m.rules = append(m.rules, r)
}

// LoadRules deserializes all the rules in the given Reader and adds them to this Manager. If there
// is any problem deserializing, no rules are added.
func (m *Manager) LoadRules(r io.Reader) error {
	dec := json.NewDecoder(r)
	var rules []*Rule
	if err := dec.Decode(&rules); err != nil {
		return err
	}
	for _, rule := range rules {
		m.AddRule(rule)
	}
	return nil
}

// DumpRules serializes all rules in this manager to the given writer.
func (m *Manager) DumpRules(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(m.rules)
}

// Evaluate runs the given transaction over all rules in the manager and applies the specified
// category when a single rule matches. The returned bool will be true if the Transaction was
// modified. A non-nil error will be returned if multiple rules matched the given Transaction.
func (m *Manager) Evaluate(t *register.Transaction) (bool, error) {
	var matched []*Rule
	for _, r := range m.rules {
		if r.Evaluate(t) {
			matched = append(matched, r)
		}
	}
	if len(matched) == 0 {
		return false, nil
	}
	if len(matched) > 1 {
		var rules []string
		for _, r := range matched {
			rules = append(rules, r.Name)
		}
		return false, fmt.Errorf("transaction %s matched multiple rules: %s", t.ID, strings.Join(rules, ", "))
	}
	if len(t.Category) > 0 {
		return false, nil
	}
	t.Category = append(t.Category, &register.Category{Name: matched[0].Category, Amount: t.Amount})
	return true, nil
}
