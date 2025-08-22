package cli

import (
	"flag"
	"fmt"
	"github.com/matrix-go/bitcoin/core"
	"os"
	"strconv"
)

type Cli struct {
	bc *core.Blockchain
}

func NewCli(bc *core.Blockchain) *Cli {
	return &Cli{bc}
}

func (c *Cli) Run() error {
	c.validateArgs()

	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "data to add")

	switch os.Args[1] {
	case "addBlock":
		if err := addBlockCmd.Parse(os.Args[2:]); err != nil {
			return err
		}
	case "printChain":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			return err
		}
	default:
		c.printUsage()
		os.Exit(1)
	}
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		c.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		c.printChain()
	}

	return nil
}

func (c *Cli) validateArgs() {
	if len(os.Args) < 2 {
		c.printUsage()
		os.Exit(1)
	}
}

func (c *Cli) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}

func (c *Cli) addBlock(data string) {
	c.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (c *Cli) printChain() {
	bci := c.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := core.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
