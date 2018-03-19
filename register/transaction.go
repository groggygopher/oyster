package register

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Transaction is a single financial transaction mapped to zero or more Categories.
type Transaction struct {
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Amount      float64     `json:"amount"`
	Date        *time.Time  `json:"date"`
	Category    []*Category `json:"categories"`
}

// String returns a quick representation of this Transaction.
func (t *Transaction) String() string {
	return fmt.Sprintf("%s,%s,$%f", t.Date, t.Description, t.Amount)
}

// ReadAllTransactions imports all transactions from a CSV file until EOF. An error will be
// returned if any error is encountered while reading or parsing.
func ReadAllTransactions(r io.Reader) ([]*Transaction, error) {
	reader := csv.NewReader(r)
	var trans []*Transaction
	// Drop header row.
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("csv.Reader.Read: %v", err)
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv.Reader.Read: %v", err)
		}
		date, err := time.Parse("1/2/2006", record[0])
		if err != nil {
			return nil, fmt.Errorf("time.Parse(%s): %v", record[0], err)
		}
		desc := record[2]
		amountStr := strings.TrimSpace(record[3] + record[4])
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return nil, fmt.Errorf("ParseFloat(%s): %v", amountStr, err)
		}
		t := &Transaction{
			ID:          fmt.Sprintf("TRANS-%s-%s-%f", date, desc, amount),
			Description: desc,
			Amount:      amount,
			Date:        &date,
		}
		trans = append(trans, t)
	}
	return trans, nil
}
