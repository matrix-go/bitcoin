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

func (b *Block) PreviousHash() [32]byte {
	return b.previousHash
}

func (b *Block) Nonce() int {
	return b.nonce
}

func (b *Block) Transactions() []*Transaction {
	return b.transactions
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

func (b *Block) UnmarshalJSON(data []byte) error {
	var val struct {
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	b.nonce = val.Nonce
	previousHash, _ := hex.DecodeString(val.PreviousHash)
	b.previousHash = [32]byte(previousHash[:32])
	b.timestamp = val.Timestamp
	b.transactions = val.Transactions
	return nil
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
