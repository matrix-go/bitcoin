package main

import (
	"github.com/matrix-go/bitcoin/core"
)

func main() {

	blockchainAddress := "miner"
	bc := core.NewBlockchain(blockchainAddress)
	bc.Print()

	bc.AddTransaction("A", "B", 10)
	bc.Mining()
	bc.Print()

	bc.AddTransaction("C", "D", 20)
	bc.AddTransaction("X", "Y", 30)
	bc.Mining()
	bc.Print()
}
