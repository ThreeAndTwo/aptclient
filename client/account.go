package client

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/mr-tron/base58"
	"github.com/threeandtwo/aptclient/key_manager"
	"github.com/threeandtwo/aptclient/types"
	"golang.org/x/crypto/sha3"
	"strings"
)

type aptAccount struct {
	key     string
	authKey string
	prvKey  ed25519.PrivateKey
	keyTy   types.KeyTy
}

func NewAptAccount(key, authKey string) *aptAccount {
	account := &aptAccount{}
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

func (a *aptAccount) prvKey2Account(prvKey string) (*types.AptAccount, error) {
	res, err := base58.Decode(prvKey)
	if err != nil {
		return nil, err
	}

	privateKey := ed25519.NewKeyFromSeed(res[:32])
	a.prvKey = privateKey

	return &types.AptAccount{
		Address:    a.address(),
		PublicKey:  a.publicKey(),
		PrivateKey: privateKey,
		AuthKey:    a.authKey,
	}, nil
}

func (a *aptAccount) genPrvKey() (*types.AptAccount, error) {
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

func (a *aptAccount) mnemonic2Account(index int) (*types.AptAccount, error) {
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

func (a *aptAccount) AccountFromRandomKey() (*types.AptAccount, error) {
	if a.keyTy != types.NoneTy {
		return nil, types.ErrNotNoneTy
	}
	return a.genPrvKey()
}

func (a *aptAccount) AccountFromPrivateKey() (*types.AptAccount, error) {
	if a.keyTy != types.PrivateTy {
		return nil, types.ErrNotPrivateKeyTy
	}
	return a.prvKey2Account(a.key)
}

func (a *aptAccount) AccountFromMnemonic(index int) (*types.AptAccount, error) {
	if a.keyTy != types.MnemonicTy {
		return nil, types.ErrNotMnemonicTy
	}

	if index < 0 {
		return nil, types.ErrMnemonicIndex
	}
	return a.mnemonic2Account(index)
}

func pubKeyBytes(prvKey ed25519.PrivateKey) []byte {
	return prvKey.Public().(ed25519.PublicKey)
}

func PrivateKey2Str(prvKey ed25519.PrivateKey) string {
	return base58.Encode(prvKey)
}

func (a *aptAccount) publicKey() string {
	return fmt.Sprint("0x", hex.EncodeToString(pubKeyBytes(a.prvKey)))
}

func (a *aptAccount) address() string {
	hasher := sha3.New256()

	hasher.Write(pubKeyBytes(a.prvKey))
	hasher.Write([]byte("\x00"))

	if a.authKey == "" {
		a.authKey = fmt.Sprint("0x", hex.EncodeToString(hasher.Sum(nil)))
	}
	return fmt.Sprint("0x", hex.EncodeToString(hasher.Sum(nil)))
}

func (a *aptAccount) Sign(msg []byte) []byte {
	return ed25519.Sign(a.prvKey, msg)
}

func (a *aptAccount) Verify(sig, msg []byte) bool {
	return ed25519.Verify(pubKeyBytes(a.prvKey), msg, sig)
}