package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    string
}

func NewWallet() *Wallet {
	w := new(Wallet)
	// 1. Create ecdsa private key (32 bytes) and public key (64 bytes)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &privateKey.PublicKey
	// 2. Perform sha256 hashing on public key
	hasher := sha256.New()
	hasher.Write(w.publicKey.X.Bytes())
	hasher.Write(w.publicKey.Y.Bytes())
	digest := hasher.Sum(nil)

	// 3. Perform RIPEMD-160 hashing on result of sha256 (20 bytes)
	md := ripemd160.New()
	md.Write(digest)
	digest = md.Sum(nil)

	// 4. Add version byte in front of RIPEMD-160 hash (0x00 for main network)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest)
	// 5. Perform sha256 for RIPEMD-160 result
	hasher = sha256.New()
	hasher.Write(vd4)
	digest = hasher.Sum(nil)
	// 6. Perform sha256 for previous step result
	hasher = sha256.New()
	hasher.Write(digest)
	digest = hasher.Sum(nil)
	// 7. Take first 4 bytes for checksum
	checksum := digest[:4]
	// 8. Add the 4 checksum to the end of extend RIPEMD-160 result with version byte
	dc8 := make([]byte, 25)
	copy(dc8, vd4)
	copy(dc8[21:], checksum)
	// 9. Convert the result from byte to base58 string
	w.address = base58.Encode(dc8)
	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%064x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) Address() string {
	return w.address
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
		Address    string `json:"address"`
	}{
		PrivateKey: w.PrivateKeyStr(),
		PublicKey:  w.PublicKeyStr(),
		Address:    w.Address(),
	})
}
