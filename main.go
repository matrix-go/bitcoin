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

	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
	fmt.Println(w.Address())
}
