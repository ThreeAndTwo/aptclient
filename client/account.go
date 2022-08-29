package client

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/mr-tron/base58"
	"github.com/threeandtwo/aptclient/hexutil"
	"github.com/threeandtwo/aptclient/key_manager"
	"github.com/threeandtwo/aptclient/types"
	"golang.org/x/crypto/sha3"
)

type AptAccount struct {
	key     string
	authKey string
	prvKey  ed25519.PrivateKey
	keyTy   types.KeyTy
}

func NewAptAccount(key, authKey string) *AptAccount {
	account := &AptAccount{}
	if IsMnemonic(key) {
		account.keyTy = types.MnemonicTy
	} else if key != "" {
		account.keyTy = types.PrivateTy
	} else {
		account.keyTy = types.NoneTy
	}

	account.key = key
	account.authKey = authKey
	return account
}

func IsMnemonic(words string) bool {
	l := len(strings.Split(words, " "))
	return l == 12 || l == 15 || l == 18 || l == 21 || l == 24
}

func (a *AptAccount) base58PrvKey2Account(prvKey string) (ed25519.PrivateKey, error) {
	res, err := base58.Decode(prvKey)
	if err != nil {
		return nil, err
	}
	return ed25519.NewKeyFromSeed(res[:32]), nil
}

func (a *AptAccount) hexPrvKey2Account(prvKey string) (ed25519.PrivateKey, error) {
	res, err := hexutil.Decode(prvKey)
	if err != nil {
		return nil, err
	}
	return ed25519.NewKeyFromSeed(res[:32]), nil
}

func (a *AptAccount) prvKey2Account(prvKey string) (*types.AptAccount, error) {
	var privateKey ed25519.PrivateKey
	var err error

	if strings.HasPrefix("0x", prvKey) || strings.HasPrefix("0X", prvKey) {
		privateKey, err = a.hexPrvKey2Account(prvKey)
	} else {
		privateKey, err = a.base58PrvKey2Account(prvKey)
	}
	if err != nil {
		return nil, err
	}

	a.prvKey = privateKey
	return &types.AptAccount{
		Address:    a.address(),
		PublicKey:  a.publicKey(),
		PrivateKey: privateKey,
		AuthKey:    a.authKey,
	}, nil
}

func (a *AptAccount) genPrvKey() (*types.AptAccount, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	a.prvKey = privateKey

	return &types.AptAccount{
		Address:    a.address(),
		PublicKey:  a.publicKey(),
		PrivateKey: privateKey,
		AuthKey:    a.authKey,
	}, nil
}

func (a *AptAccount) mnemonic2Account(index int) (*types.AptAccount, error) {
	km, err := key_manager.NewKeyManagerWithMnemonic(256, "", a.key)
	if err != nil {
		return nil, types.ErrMnemonicCount
	}
	key, err := km.GetKey(key_manager.PurposeBIP44, key_manager.CoinTypeAPT, 0, 0, uint32(index))
	if err != nil {
		return nil, types.ErrMnemonicCount
	}

	_, _, prvKey := key.EncodeEth()
	seed, err := hex.DecodeString(prvKey)
	if err != nil {
		return nil, err
	}

	privateKey := ed25519.NewKeyFromSeed(seed[:32])
	a.prvKey = privateKey
	return &types.AptAccount{
		Address:    a.address(),
		PublicKey:  a.publicKey(),
		PrivateKey: privateKey,
		AuthKey:    a.authKey,
	}, nil
}

func (a *AptAccount) AccountFromRandomKey() (*types.AptAccount, error) {
	if a.keyTy != types.NoneTy {
		return nil, types.ErrNotNoneTy
	}
	return a.genPrvKey()
}

func (a *AptAccount) AccountFromPrivateKey() (*types.AptAccount, error) {
	if a.keyTy != types.PrivateTy {
		return nil, types.ErrNotPrivateKeyTy
	}
	return a.prvKey2Account(a.key)
}

func (a *AptAccount) AccountFromMnemonic(index int) (*types.AptAccount, error) {
	if a.keyTy != types.MnemonicTy {
		return nil, types.ErrNotMnemonicTy
	}

	if index < 0 {
		return nil, types.ErrMnemonicIndex
	}
	return a.mnemonic2Account(index)
}

func (a *AptAccount) GetAptAccount(index int) (*types.AptAccount, error) {
	switch a.keyTy {
	case types.MnemonicTy:
		return a.AccountFromMnemonic(index)
	case types.PrivateTy:
		return a.AccountFromPrivateKey()
	case types.NoneTy:
		return a.AccountFromRandomKey()
	default:
		return nil, types.ErrNotNoneTy
	}
}

func pubKeyBytes(prvKey ed25519.PrivateKey) []byte {
	return prvKey.Public().(ed25519.PublicKey)
}

func PrivateKey2Str(prvKey ed25519.PrivateKey) string {
	return base58.Encode(prvKey)
}

func (a *AptAccount) publicKey() string {
	return fmt.Sprint("0x", hex.EncodeToString(pubKeyBytes(a.prvKey)))
}

func (a *AptAccount) address() string {
	hasher := sha3.New256()

	hasher.Write(pubKeyBytes(a.prvKey))
	hasher.Write([]byte("\x00"))

	if a.authKey == "" {
		a.authKey = fmt.Sprint("0x", hex.EncodeToString(hasher.Sum(nil)))
	}
	return fmt.Sprint("0x", hex.EncodeToString(hasher.Sum(nil)))
}

func (a *AptAccount) Sign(msg []byte) []byte {
	return ed25519.Sign(a.prvKey, msg)
}

func (a *AptAccount) Verify(sig, msg []byte) bool {
	return ed25519.Verify(pubKeyBytes(a.prvKey), msg, sig)
}
