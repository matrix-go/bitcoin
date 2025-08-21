package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	sender    string
	recipient string
	value     uint64
}

func NewTransaction(sender string, recipient string, value uint64) *Transaction {
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
		Value     uint64 `json:"value"`
	}{
		Sender:    t.sender,
		Recipient: t.recipient,
		Value:     t.value,
	})
}

func (t *Transaction) Print() {
	fmt.Printf(`%s
sender         %s
recipient      %s
value   	   %d
`, strings.Repeat("-", 20), t.sender, t.recipient, t.value)
}
