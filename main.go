package main

import (
	"github.com/matrix-go/bitcoin/core"
)

func main() {
	bc := core.NewBlockchain()
	bc.Print()

	bc.AddTransaction("A", "B", 10)
	previousHash := bc.LastBlock().Hash()
	nonce := bc.ProofOfWork()
	bc.CreateBlock(nonce, previousHash)
	bc.Print()

	bc.AddTransaction("C", "D", 20)
	bc.AddTransaction("X", "Y", 30)
	previousHash = bc.LastBlock().Hash()
	nonce = bc.ProofOfWork()
	bc.CreateBlock(nonce, previousHash)
	bc.Print()
}
