package core

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/matrix-go/bitcoin/utils"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	KMiningDifficulty = 3
	KMiningSender     = "COINBASE"
	KMiningReward     = 20
	KMiningTimerSec   = 20

	KBlockchainPortStart      = 5000
	KBlockchainPortEnd        = 5003
	KNeighborIPRangeStart     = 0
	KNeighborIPRangeEnd       = 1
	KBlockNeighborSyncTimeSec = 20
)

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
	miner           string
	port            int
	mux             sync.Mutex

	neighbors   []string
	muxNeighbor sync.Mutex
}

func NewBlockchain(miner string, port int) *Blockchain {
	bc := new(Blockchain)
	bc.transactionPool = make([]*Transaction, 0)
	bc.chain = make([]*Block, 0)
	bc.miner = miner
	bc.port = port

	// TODO: genesis block
	b := &Block{}
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) Run() {
	bc.SyncNeighbors()
}

func (bc *Blockchain) ResolveConflicts() bool {
	var longestChain []*Block
	var maxLength = len(bc.chain)
	for _, nb := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/blockchain", nb)
		resp, _ := http.Get(endpoint)
		if resp.StatusCode == http.StatusOK {
			var blockchain Blockchain
			_ = json.NewDecoder(resp.Body).Decode(&blockchain)
			chain := blockchain.Chain()
			if len(chain) > maxLength {
				longestChain = chain
				maxLength = len(longestChain)
			}
		}
	}
	if longestChain != nil {
		bc.chain = longestChain
		log.Printf("Resolve conflicts with replace chain - chain length: %d\n", len(longestChain))
		return true
	}
	log.Printf("Resolve conflicts not replace chain\n")
	return false
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)

	// TODO: transaction
	bc.transactionPool = make([]*Transaction, 0)
	for _, neighbor := range bc.neighbors {
		endpoint := fmt.Sprintf("http://%s/transactions", neighbor)
		req, _ := http.NewRequest("DELETE", endpoint, nil)
		req.Header.Add("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		log.Printf("%v\n", resp)
	}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddTransaction(
	sender string,
	recipient string,
	value int64,

	senderPublicKey *ecdsa.PublicKey,
	sig *utils.Signature,
) bool {

	tx := NewTransaction(sender, recipient, value)
	if sender == KMiningSender {
		// TODO: coinbase tx has no signature
		bc.transactionPool = append(bc.transactionPool, tx)
		return true
	}
	if bc.VerifyTransaction(senderPublicKey, sig, tx) {
		//if bc.CalculateTotalAmount(sender) < value {
		//	log.Println("ERROR: Not Enough balance for the sender", sender)
		//	return false
		//}
		bc.transactionPool = append(bc.transactionPool, tx)
		return true
	}
	log.Println("ERROR: AddTransaction Failed")
	return false
}

func (bc *Blockchain) VerifyTransaction(
	senderPublicKey *ecdsa.PublicKey,
	sig *utils.Signature,
	tx *Transaction,
) bool {
	m, _ := json.Marshal(tx)
	hash := sha256.Sum256(m)
	return ecdsa.Verify(senderPublicKey, hash[:], sig.R, sig.S)
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

func (bc *Blockchain) GetTransactionPools() []*Transaction {
	return bc.transactionPool
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
	bc.mux.Lock()
	defer bc.mux.Unlock()

	// TODO: in fact empty tx will also mining a new block
	if len(bc.transactionPool) == 0 {
		return false
	}

	bc.AddTransaction(KMiningSender, bc.miner, KMiningReward, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	return true
}

func (bc *Blockchain) StartMining() {
	_ = bc.Mining()
	_ = time.AfterFunc(KMiningTimerSec*time.Second, bc.StartMining)
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) int64 {
	var total int64
	for _, block := range bc.chain {
		for _, tx := range block.transactions {
			if tx.recipient == blockchainAddress {
				total += tx.value
			}
			if tx.sender == blockchainAddress {
				total -= tx.value
			}
		}
	}
	return total
}

func (bc *Blockchain) SetNeighbors() {
	bc.neighbors = utils.FindNeighbors(
		utils.GetHost(), bc.port,
		KNeighborIPRangeStart, KNeighborIPRangeEnd,
		KBlockchainPortStart, KBlockchainPortEnd,
	)
}

func (bc *Blockchain) SyncNeighbors() {
	bc.muxNeighbor.Lock()
	defer bc.muxNeighbor.Unlock()
	bc.SetNeighbors()
}

func (bc *Blockchain) StartSyncNeighbors() {
	bc.SyncNeighbors()
	_ = time.AfterFunc(time.Second*KBlockNeighborSyncTimeSec, bc.StartSyncNeighbors)
}

func (bc *Blockchain) ValidChain(chain []*Block) bool {
	prevBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		block := chain[currentIndex]
		if block.previousHash != prevBlock.Hash() {
			return false
		}
		if !bc.VerifyProof(block.nonce, block.previousHash, block.transactions, KMiningDifficulty) {
			return false
		}
		prevBlock = block
		currentIndex++
	}
	return true
}

func (bc *Blockchain) Chain() []*Block {
	return bc.chain
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"blocks"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) UnmarshalJSON(data []byte) error {
	var val struct {
		Blocks []*Block `json:"blocks"`
	}
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	bc.chain = val.Blocks
	return nil
}

func (bc *Blockchain) Print() {
	fmt.Printf("%s Blockchain %s\n", strings.Repeat("*", 30), strings.Repeat("*", 30))
	for _, b := range bc.chain {
		b.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 72))
}
