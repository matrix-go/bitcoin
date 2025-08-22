package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	sender    string
	recipient string
	value     int64
}

func NewTransaction(sender string, recipient string, value int64) *Transaction {
	tx := new(Transaction)
	tx.sender = sender
	tx.recipient = recipient
	tx.value = value
	return tx
}

// MarshalJSON implement json.Marshaler
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string `json:"sender"`
		Recipient string `json:"recipient"`
		Value     int64  `json:"value"`
	}{
		Sender:    t.sender,
		Recipient: t.recipient,
		Value:     t.value,
	})
}

func (t *Transaction) UnmarshalJSON(b []byte) error {
	var val struct {
		Sender    string `json:"sender"`
		Recipient string `json:"recipient"`
		Value     int64  `json:"value"`
	}
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}
	t.sender = val.Sender
	t.recipient = val.Recipient
	t.value = val.Value
	return nil
}

func (t *Transaction) Print() {
	fmt.Printf(`%s
sender         %s
recipient      %s
value   	   %d
`, strings.Repeat("-", 20), t.sender, t.recipient, t.value)
}
