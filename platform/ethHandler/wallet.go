package ethHandler

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    common.Address
}

// TODO: Convert to LRU cache with a fixed capacity
var walletCache = map[string]*Wallet{}

func GetWallet(privateKeyHex string) (*Wallet, error) {
	cachedWallet, cachedWalletFound := walletCache[privateKeyHex]
	if cachedWalletFound {
		return cachedWallet, nil
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	wallet := &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKeyECDSA,
		Address:    address,
	}

	walletCache[privateKeyHex] = wallet
	return wallet, nil
}

func (w *Wallet) PrivateKeyHex() string {
	privateKeyBytes := crypto.FromECDSA(w.PrivateKey)
	return hexutil.Encode(privateKeyBytes)[2:]
}

func (w *Wallet) PublicKeyHex() string {
	publicKeyBytes := crypto.FromECDSAPub(w.PublicKey)
	return hexutil.Encode(publicKeyBytes)[4:]
}

func (w *Wallet) String() string {
	return w.Address.Hex()
}
