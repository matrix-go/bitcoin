package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"github.com/matrix-go/bitcoin/utils"
)

type Transaction struct {
	senderPrivateKey *ecdsa.PrivateKey
	senderPublicKey  *ecdsa.PublicKey
	senderAddress    string
	recipientAddress string
	value            int64
}

func NewTransaction(
	senderPrivateKey *ecdsa.PrivateKey,
	senderPublicKey *ecdsa.PublicKey,
	senderAddress string,
	receiptAddress string,
	value int64,
) *Transaction {
	return &Transaction{
		senderPrivateKey: senderPrivateKey,
		senderPublicKey:  senderPublicKey,
		senderAddress:    senderAddress,
		recipientAddress: receiptAddress,
		value:            value,
	}
}

func (tx *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(tx)
	h := sha256.Sum256(m)
	r, s, _ := ecdsa.Sign(rand.Reader, tx.senderPrivateKey, h[:])
	return &utils.Signature{
		R: r,
		S: s,
	}
}

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string `json:"sender"`
		Recipient string `json:"recipient"`
		Value     int64  `json:"value"`
	}{
		Sender:    tx.senderAddress,
		Recipient: tx.recipientAddress,
		Value:     tx.value,
	})
}
