package utils

import (
	"encoding/hex"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (sig *Signature) String() string {
	return hex.EncodeToString(sig.R.Bytes()) + hex.EncodeToString(sig.S.Bytes())
}
