package core

import (
	"go.etcd.io/bbolt"
	"log"
)

var (
	dbFile       = "blockchain.db"
	blocksBucket = "blocks"
)

type Blockchain struct {
	tip []byte
	db  *bbolt.DB
}

func NewBlockchain() *Blockchain {

	var tip []byte
	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bbolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			genesis := NewGenesisBlock()
			if bucket, err = tx.CreateBucket([]byte(blocksBucket)); err != nil {
				return err
			}
			if err = bucket.Put(genesis.Hash, genesis.Serialize()); err != nil {
				return err
			}
			if err = bucket.Put([]byte("l"), genesis.Hash); err != nil {
				return err
			}
			tip = genesis.Hash
		} else {
			// l -> current block hash
			tip = bucket.Get([]byte("l"))
		}
		return nil
	})
	bc := &Blockchain{tip, db}
	return bc
}

func (bc *Blockchain) AddBlock(data string) error {
	var lastHash []byte
	err := bc.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		return err
	}
	newBlock := NewBlock(data, lastHash)
	return bc.db.Update(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(blocksBucket))
		if err = b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			return err
		}
		if err = b.Put([]byte("l"), newBlock.Hash); err != nil {
			return err
		}
		bc.tip = newBlock.Hash
		return nil
	})
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bbolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

func (bc *Blockchain) Close() {
	bc.db.Close()
}

func (it *BlockchainIterator) Next() *Block {
	var block *Block
	err := it.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		block = DeserializeBlock(b.Get(it.currentHash))
		return nil
	})
	if err != nil {
		log.Printf("get block err: %s", err)
		return nil
	}
	it.currentHash = block.PrevBlockHash
	return block
}
