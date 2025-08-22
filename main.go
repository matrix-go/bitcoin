package main

import (
	"github.com/matrix-go/bitcoin/cli"
	"github.com/matrix-go/bitcoin/core"
	"log"
)

func main() {
	bc := core.NewBlockchain()
	defer bc.Close()

	c := cli.NewCli(bc)
	log.Fatal(c.Run())
}
