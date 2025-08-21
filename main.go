package main

import (
	"fmt"
	"github.com/matrix-go/bitcoin/core"
	"github.com/matrix-go/bitcoin/wallet"
)

func main() {

	// Miner
	walletM := wallet.NewWallet()

	// A and B
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	// tx
	var amount int64 = 100
	tx := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.Address(), walletB.Address(), amount)

	// bc
	bc := core.NewBlockchain(walletM.Address())

	success := bc.AddTransaction(walletA.Address(), walletB.Address(), amount, walletA.PublicKey(), tx.GenerateSignature())
	fmt.Println("is success", success)
	bc.Mining()
	bc.Print()

	fmt.Printf("Miner: %d\n", bc.CalculateTotalAmount(walletM.Address()))
	fmt.Printf("A: %d\n", bc.CalculateTotalAmount(walletA.Address()))
	fmt.Printf("B: %d\n", bc.CalculateTotalAmount(walletB.Address()))
}
