package session

import (
	"reflect"
	"testing"
	"time"

	"github.com/groggygopher/oyster/register"
)

func TestSerializeDeserialize(t *testing.T) {
	usr := &User{
		Name: "test",
		transactions: []*register.Transaction{
			&register.Transaction{
				Description: "test",
			},
		},
	}

	bs, err := usr.Serialize()
	if err != nil {
		t.Fatalf("user.Serialize: %v", err)
	}

	deser, err := DeserializeUser(bs)
	if err != nil {
		t.Fatalf("DeserializeUser: %v", err)
	}

	if got, want := deser, usr; !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestTransactions(t *testing.T) {
	var (
		today = time.Now()
		back1 = time.Now().Add(-24 * time.Hour)
		back2 = time.Now().Add(-48 * time.Hour)
		back3 = time.Now().Add(-72 * time.Hour)
	)
	var (
		trans1 = &register.Transaction{
			ID:          "trans1",
			Description: "trans1",
			Amount:      1.23,
			Date:        &today,
		}
		trans2 = &register.Transaction{
			ID:          "trans2",
			Description: "trans2",
			Amount:      2.23,
			Date:        &back1,
		}
		trans3 = &register.Transaction{
			ID:          "trans3",
			Description: "trans3",
			Amount:      3.23,
			Date:        &back2,
		}
		trans4 = &register.Transaction{
			ID:          "trans4",
			Description: "trans4",
			Amount:      4.23,
			Date:        &back3,
		}
	)

	tests := []struct {
		label        string
		firstImport  []*register.Transaction
		secondImport []*register.Transaction
		wantTrans    []*register.Transaction
		firstCount   int
		secondCount  int
	}{
		{
			label:        "empty",
			firstImport:  []*register.Transaction{},
			secondImport: []*register.Transaction{},
		},
		{
			label: "all first",
			firstImport: []*register.Transaction{
				trans1,
				trans2,
				trans3,
				trans4,
			},
			secondImport: []*register.Transaction{},
			wantTrans: []*register.Transaction{
				trans1,
				trans2,
				trans3,
				trans4,
			},
			firstCount:  4,
			secondCount: 0,
		},
		{
			label: "two, not overlapping",
			firstImport: []*register.Transaction{
				trans1,
			},
			secondImport: []*register.Transaction{
				trans4,
			},
			wantTrans: []*register.Transaction{
				trans1,
				trans4,
			},
			firstCount:  1,
			secondCount: 1,
		},
		{
			label: "all, not overlapping",
			firstImport: []*register.Transaction{
				trans1,
				trans2,
			},
			secondImport: []*register.Transaction{
				trans3,
				trans4,
			},
			wantTrans: []*register.Transaction{
				trans1,
				trans2,
				trans3,
				trans4,
			},
			firstCount:  2,
			secondCount: 2,
		},
		{
			label: "one overlapping",
			firstImport: []*register.Transaction{
				trans1,
				trans2,
				trans3,
			},
			secondImport: []*register.Transaction{
				trans3,
				trans4,
			},
			wantTrans: []*register.Transaction{
				trans1,
				trans2,
				trans3,
				trans4,
			},
			firstCount:  3,
			secondCount: 1,
		},
		{
			label: "all overlapping",
			firstImport: []*register.Transaction{
				trans1,
				trans2,
				trans3,
				trans4,
			},
			secondImport: []*register.Transaction{
				trans1,
				trans2,
				trans3,
				trans4,
			},
			wantTrans: []*register.Transaction{
				trans1,
				trans2,
				trans3,
				trans4,
			},
			firstCount: 4,
		},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			usr := &User{}
			if got, want := usr.ImportTransactions(test.firstImport), test.firstCount; got != want {
				t.Errorf("first import count mismatch: got: %d, want: %d", got, want)
			}
			if got, want := usr.ImportTransactions(test.secondImport), test.secondCount; got != want {
				t.Errorf("second import count mismatch: got: %d, want: %d", got, want)
			}
			if got, want := usr.Transactions(), test.wantTrans; !reflect.DeepEqual(got, want) {
				t.Errorf("transactions: got: %v, want: %v", got, want)
			}
		})
	}
}
