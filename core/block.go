package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	// TODO: transaction
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.nonce = nonce
	b.previousHash = previousHash
	b.timestamp = time.Now().UnixNano()
	b.transactions = transactions
	return b
}

func (b *Block) Hash() [32]byte {
	// TODO: encoder
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

// MarshalJSON implement json.Marshaler
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PreviousHash: hex.EncodeToString(b.previousHash[:]),
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

func (b *Block) Print() {
	fmt.Printf(`%s Block %s
timestamp      %d
nonce          %d
previousHash   %x
`, strings.Repeat("=", 30), strings.Repeat("=", 30), b.timestamp, b.nonce, b.previousHash)
	if len(b.transactions) == 0 {
		fmt.Println("transactions   []")
	}
	for _, tx := range b.transactions {
		tx.Print()
	}
}
