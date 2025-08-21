package core

import (
	"fmt"
	"strings"
)

const (
	KMiningDifficulty = 3
	KMiningSender     = "COINBASE"
	KMiningReward     = 20
)

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	bc := new(Blockchain)
	bc.transactionPool = make([]*Transaction, 0)
	bc.chain = make([]*Block, 0)
	bc.blockchainAddress = blockchainAddress

	// TODO: genesis block
	b := &Block{}
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)

	// TODO: transaction
	bc.transactionPool = make([]*Transaction, 0)
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value uint64) {
	tx := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, tx)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, tx := range bc.transactionPool {
		transactions = append(
			transactions,
			NewTransaction(tx.sender, tx.recipient, tx.value),
		)
	}
	return transactions
}

func (bc *Blockchain) VerifyProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := &Block{
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
	guessHashString := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashString[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()

	nonce := 0
	for !bc.VerifyProof(nonce, previousHash, transactions, KMiningDifficulty) {
		nonce++
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(KMiningSender, bc.blockchainAddress, KMiningReward)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	return true
}

func (bc *Blockchain) Print() {
	fmt.Printf("%s Blockchain %s\n", strings.Repeat("*", 30), strings.Repeat("*", 30))
	for _, b := range bc.chain {
		b.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 72))
}
