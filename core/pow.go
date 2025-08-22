package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

func (p *ProofOfWork) PrepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			p.block.PrevBlockHash,
			p.block.Data,
			IntToHex(p.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(nonce),
		},
		[]byte{},
	)
	return data
}

func (p *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64 = 0

	fmt.Printf("Mining the block containing \"%s\"\n", p.block.Data)
	for nonce < math.MaxInt64 {
		data := p.PrepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(p.target) == -1 {
			break
		}
		nonce++
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

func (p *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := p.PrepareData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(p.target) == -1
	return isValid
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
