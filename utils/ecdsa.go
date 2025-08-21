package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (sig *Signature) String() string {
	return fmt.Sprintf("%064x%064x", sig.R.Bytes(), sig.S.Bytes())
}

func String2BigIntTuple(str string) (big.Int, big.Int) {
	bx, _ := hex.DecodeString(str[:64])
	by, _ := hex.DecodeString(str[64:])

	var x big.Int
	x.SetBytes(bx)
	var y big.Int
	y.SetBytes(by)
	return x, y
}

func PublicKeyFromString(str string) (*ecdsa.PublicKey, error) {
	x, y := String2BigIntTuple(str)
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}, nil
}

func PrivateKeyFromString(str string, publicKey ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	b, _ := hex.DecodeString(str)
	var d big.Int
	d.SetBytes(b)
	return &ecdsa.PrivateKey{
		PublicKey: publicKey,
		D:         &d,
	}, nil
}

func SignatureFromString(str string) (*Signature, error) {
	r, s := String2BigIntTuple(str)
	return &Signature{R: &r, S: &s}, nil
}
