package main

import (
	"fmt"
	"github.com/matrix-go/bitcoin/wallet"
)

func main() {

	//blockchainAddress := "miner"
	//bc := core.NewBlockchain(blockchainAddress)
	//bc.Print()
	//
	//bc.AddTransaction("A", "B", 10)
	//bc.Mining()
	//bc.Print()
	//
	//bc.AddTransaction("C", "D", 20)
	//bc.AddTransaction("X", "Y", 30)
	//bc.Mining()
	//bc.Print()
	//
	//fmt.Printf("C: %d\n", bc.CalculateTotalAmount("C"))
	//fmt.Printf("D: %d\n", bc.CalculateTotalAmount("D"))

	a := wallet.NewWallet()
	fmt.Println(a.PrivateKeyStr())
	fmt.Println(a.PublicKeyStr())
	fmt.Println(a.Address())

	b := wallet.NewWallet()
	tx := wallet.NewTransaction(a.PrivateKey(), a.PublicKey(), a.Address(), b.Address(), 10)
	fmt.Printf("tx ==> %v\n", tx)
	sig := tx.GenerateSignature()
	fmt.Printf("sig ==> %v\n", sig)
}
