package client

import (
	"errors"
	"fmt"
	"github.com/threeandtwo/aptclient/hexutil"
	"github.com/threeandtwo/aptclient/types"
	"os"
	"testing"
)

func TestNewAptAccount(t *testing.T) {
	var tests = []struct {
		desc    string
		key     string
		authKey string
		index   int
		keyTy   types.KeyTy
	}{
		{
			desc:    "12L words",
			key:     os.Getenv("KEY"),
			authKey: os.Getenv("AUTH_KEY"),
			index:   0,
			keyTy:   types.MnemonicTy,
		},
		{
			desc:    "null key",
			key:     "",
			authKey: "",
			index:   0,
			keyTy:   types.NoneTy,
		},
		{
			desc:    "correct privateKey",
			key:     os.Getenv("KEY"),
			authKey: os.Getenv("AUTH_KEY"),
			index:   0,
			keyTy:   types.PrivateTy,
		},
		{
			desc:    "error privateKey",
			key:     os.Getenv("KEY"),
			authKey: os.Getenv("AUTH_KEY"),
			index:   0,
			keyTy:   types.PrivateTy,
		},
		{
			desc:    "error mnemonic",
			key:     os.Getenv("KEY"),
			authKey: os.Getenv("AUTH_KEY"),
			index:   0,
			keyTy:   types.MnemonicTy,
		},
		{
			desc:    "index < 0 for mnemonic",
			key:     os.Getenv("KEY"),
			authKey: os.Getenv("AUTH_KEY"),
			index:   -3,
			keyTy:   types.MnemonicTy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			na := NewAptAccount(tt.key, tt.authKey)

			if na.keyTy != tt.keyTy {
				t.Logf("keyTy mismatched for %s \n", tt.desc)
				return
			}

			var err error
			switch na.keyTy {
			case types.MnemonicTy:
				err = accountByMnemonic(na, tt.index)
			case types.PrivateTy:
				err = accountByPrivateKey(na)
			case types.NoneTy:
				err = accountByRandomKey(na)
			default:
				t.Logf("%d keyTy unsupport for %s", tt.keyTy, tt.desc)
			}

			if err != nil {
				t.Logf("%s error for %s", err.Error(), tt.desc)
			}
		})
	}
}

func TestAptAccount_GetAptAccount(t *testing.T) {
	var tests = []struct {
		name    string
		key     string
		authKey string
		index   int
		keyTy   types.KeyTy
	}{
		{
			name:    "12L words",
			key:     os.Getenv("KEY"),
			authKey: os.Getenv("AUTH_KEY"),
			index:   0,
			keyTy:   types.MnemonicTy,
		},
		{
			name:    "null key",
			key:     "",
			authKey: "",
			index:   0,
			keyTy:   types.NoneTy,
		},
		{
			name:    "correct privateKey",
			key:     os.Getenv("KEY"),
			authKey: os.Getenv("AUTH_KEY"),
			index:   0,
			keyTy:   types.PrivateTy,
		},
		{
			name:    "error privateKey",
			key:     os.Getenv("KEY"),
			authKey: os.Getenv("AUTH_KEY"),
			index:   0,
			keyTy:   types.PrivateTy,
		},
		{
			name:    "error mnemonic",
			key:     os.Getenv("KEY"),
			authKey: os.Getenv("AUTH_KEY"),
			index:   0,
			keyTy:   types.MnemonicTy,
		},
		{
			name:    "index < 0 for mnemonic",
			key:     os.Getenv("KEY"),
			authKey: os.Getenv("AUTH_KEY"),
			index:   -3,
			keyTy:   types.MnemonicTy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			na := NewAptAccount(tt.key, tt.authKey)

			if na.keyTy != tt.keyTy {
				t.Logf("keyTy mismatched for %s \n", tt.name)
				return
			}

			_account, err := na.GetAptAccount(tt.index)
			if err != nil {
				t.Logf("keyTy mismatched for %s \n", tt.name)
				return
			}
			printAccount(_account)
		})
	}
}

func accountByMnemonic(a *AptAccount, index int) error {
	_account, err := a.AccountFromMnemonic(index)
	if err != nil && err != types.ErrMnemonicIndex {
		return err
	}

	printAccount(_account)

	msg := "This is a sample message"
	sign := a.Sign([]byte(msg))
	fmt.Printf("✍️ signature: %s \n", hexutil.Encode(sign))

	_verify := a.Verify(sign, []byte(msg))
	if !_verify {
		return errors.New("signature verify error ❌")
	}
	fmt.Printf("signature verified ✅ \n")
	return nil
}

func printAccount(a *types.AptAccount) {
	fmt.Printf("address: %s \n", a.Address)
	fmt.Printf("privateKey: %s \n", PrivateKey2Str(a.PrivateKey))
	fmt.Printf("authKey: %s \n", a.AuthKey)
	fmt.Printf("publicKey: %s \n", a.PublicKey)
}

func accountByPrivateKey(a *AptAccount) error {
	_account, err := a.AccountFromPrivateKey()
	if err != nil {
		return err
	}

	printAccount(_account)

	msg := "This is a sample message"
	sign := a.Sign([]byte(msg))
	fmt.Printf("✍️ signature: %s \n", hexutil.Encode(sign))

	_verify := a.Verify(sign, []byte(msg))
	if !_verify {
		return errors.New("signature verify error ❌")
	}
	fmt.Printf("signature verified ✅ \n")
	return nil
}

func accountByRandomKey(a *AptAccount) error {
	_account, err := a.AccountFromRandomKey()
	if err != nil {
		return err
	}

	printAccount(_account)
	msg := "This is a sample message"
	sign := a.Sign([]byte(msg))
	fmt.Printf("✍️ signature: %s \n", hexutil.Encode(sign))

	_verify := a.Verify(sign, []byte(msg))
	if !_verify {
		return errors.New("signature verify error ❌")
	}
	fmt.Printf("signature verified ✅ \n")

	return nil
}
