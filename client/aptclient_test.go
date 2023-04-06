package client

import (
	"encoding/json"
	"fmt"
	"github.com/threeandtwo/aptclient/types"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	RPC_ADDR         string = "https://fullnode.devnet.aptoslabs.com/v1"
	MAINNET_RPC_ADDR string = "https://fullnode.mainnet.aptoslabs.com/v1"
)

func TestNewAptClient(t *testing.T) {
	tests := []struct {
		name string
		rpc  string
	}{
		{
			name: "correct rpc addr",
			rpc:  RPC_ADDR,
		},
		{
			name: "incorrect rpc addr",
			rpc:  "",
		},
		{
			name: "null faucet url",
			rpc:  RPC_ADDR + "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}
			t.Logf("rpc addr: %s", c.rpc)
		})
	}
}

func TestAptClient_NodeHealth(t *testing.T) {
	tests := []struct {
		name string
		rpc  string
		secs uint32
	}{
		{
			name: "test node health for aptos, duration is 0",
			rpc:  RPC_ADDR,
			secs: 0,
		},
		{
			name: "secs > 0",
			rpc:  RPC_ADDR,
			secs: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			healthDesc, err := c.NodeHealth(tt.secs)
			if err != nil {
				t.Logf("node health is error: %s", err.Error())
				return
			}

			t.Logf("message: %s", healthDesc)
		})
	}
}

func TestAptClient_LedgerInfo(t *testing.T) {
	tests := []struct {
		name string
		rpc  string
	}{
		{
			name: "test ledgerInfo",
			rpc:  RPC_ADDR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			info, err := c.LedgerInfo()
			if err != nil {
				t.Logf("get ledgerInfo: %s", err)
				return
			}

			b, _ := json.Marshal(info)
			t.Logf("resources: %s", string(b))
		})
	}
}

func TestAptClient_BlockByHeight(t *testing.T) {
	tests := []struct {
		name        string
		blockHeight uint64
		withTxs     types.BlockWithTxs
	}{
		{
			name:        "blockHeight <= 0",
			blockHeight: 0,
			withTxs:     types.FalseTy,
		},
		{
			name:        "block had been confirmed",
			blockHeight: 11111,
			withTxs:     types.FalseTy,
		},
		{
			name:        "block had been confirmed and withTxs is True",
			blockHeight: 11111,
			withTxs:     types.TrueTy,
		},
		{
			name:        "block will be minted",
			blockHeight: 1000000000000000,
			withTxs:     types.FalseTy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(RPC_ADDR)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			block, err := c.BlockByHeight(tt.blockHeight, tt.withTxs)
			if err != nil {
				t.Logf("get block info error: %s", err.Error())
				return
			}

			b, _ := json.Marshal(block)
			t.Logf("resources: %s", string(b))
		})
	}
}

func TestAptClient_BlockByVersion(t *testing.T) {
	tests := []struct {
		name    string
		version uint64
		withTxs types.BlockWithTxs
	}{
		{
			name:    "version == 0",
			version: 0,
			withTxs: types.TrueTy,
		},
		{
			name:    "version for block had been confirmed",
			version: 11111,
			withTxs: types.FalseTy,
		},
		{
			name:    "version for block had been confirmed and withTxs is true",
			version: 11111,
			withTxs: types.TrueTy,
		},
		{
			name:    "version for block will be minted",
			version: 1000000000000000,
			withTxs: types.TrueTy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(RPC_ADDR)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			block, err := c.BlockByVersion(tt.version, tt.withTxs)
			if err != nil {
				t.Logf("get block info error: %s", err.Error())
				return
			}

			b, _ := json.Marshal(block)
			t.Logf("resources: %s", string(b))
		})
	}
}

func TestAptClient_Account(t *testing.T) {
	tests := []struct {
		name    string
		rpc     string
		address string
	}{
		{
			name:    "test Account",
			rpc:     MAINNET_RPC_ADDR,
			address: "0x9a61496f9603d4c48d1ccb6d1aebcffefbd0efd6270e46885744c191ef2bfd59",
		},
		{
			name:    "address is null",
			rpc:     RPC_ADDR,
			address: "",
		},
		{
			name:    "address length error",
			rpc:     RPC_ADDR,
			address: "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return

			}
			info, err := c.Account(tt.address)
			if err != nil {
				t.Logf("get account Info: %s", err)
				return
			}
			t.Logf("ledgerInfo: %v", info)
		})
	}
}

func TestAptClient_Balance(t *testing.T) {
	tests := []struct {
		name    string
		rpc     string
		address string
	}{
		{
			name:    "test balance",
			rpc:     RPC_ADDR,
			address: "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
		},
		{
			name:    "address is null",
			rpc:     RPC_ADDR,
			address: "",
		},
		{
			name:    "address length error",
			rpc:     RPC_ADDR,
			address: "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return

			}
			info, err := c.GetBalance(tt.address)
			if err != nil {
				t.Logf("get balance Info error: %s", err)
				return
			}
			t.Logf("balance %sAPT for %s", info.String(), tt.address)
		})
	}
}

func TestAptClient_GetNonce(t *testing.T) {
	tests := []struct {
		name    string
		rpc     string
		address string
	}{
		{
			name:    "test balance",
			rpc:     RPC_ADDR,
			address: "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
		},
		{
			name:    "address is null",
			rpc:     RPC_ADDR,
			address: "",
		},
		{
			name:    "address length error",
			rpc:     RPC_ADDR,
			address: "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return

			}
			info, err := c.GetNonce(tt.address)
			if err != nil {
				t.Logf("get balance Info error: %s", err)
				return
			}
			t.Logf("nonce %d for %s", info, tt.address)
		})
	}
}

func TestAptClient_AccountResources(t *testing.T) {
	tests := []struct {
		name    string
		rpc     string
		address string
		version string
	}{
		{
			name:    "normal account",
			rpc:     RPC_ADDR,
			address: "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
			version: "7935073",
		},
		{
			name:    "version is null",
			rpc:     RPC_ADDR,
			address: "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
			version: "",
		},
		{
			name:    "address is null",
			rpc:     RPC_ADDR,
			address: "",
			version: "7920450",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			resources, err := c.AccountResources(tt.address, tt.version)
			if err != nil {
				t.Logf("get account resource for %s error: %s", tt.name, err.Error())
				return
			}

			b, _ := json.Marshal(resources)
			t.Logf("resources: %s", string(b))
		})
	}
}

func TestAptClient_AccountResourceByType(t *testing.T) {
	tests := []struct {
		name       string
		rpc        string
		address    string
		resourceTy string
		version    string
	}{
		{
			name:       "normal",
			rpc:        RPC_ADDR,
			address:    "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
			resourceTy: "0x1::account::Account",
			version:    "7920450",
		},
		{
			name:       "version is null",
			rpc:        RPC_ADDR,
			address:    "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
			resourceTy: "0x1::account::Account",
			version:    "",
		},
		{
			name:       "resource type is null",
			rpc:        RPC_ADDR,
			address:    "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
			resourceTy: "",
			version:    "7920450",
		},
		{
			name:       "address is null",
			rpc:        RPC_ADDR,
			address:    "",
			resourceTy: "0x1::account::Account",
			version:    "7920450",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			resource, err := c.AccountResourceByType(tt.address, tt.resourceTy, tt.version)
			if err != nil {
				t.Logf("get account resource for %s error: %s", tt.name, err.Error())
				return
			}

			b, _ := json.Marshal(resource)
			t.Logf("resource: %s", string(b))
		})
	}
}

func TestAptClient_AccountModules(t *testing.T) {
	tests := []struct {
		name    string
		rpc     string
		address string
		version string
	}{
		{
			name:    "normal",
			rpc:     RPC_ADDR,
			address: "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
			version: "7920450",
		},
		{
			name:    "address is null",
			rpc:     RPC_ADDR,
			address: "",
			version: "7920450",
		},
		{
			name:    "version is null",
			rpc:     RPC_ADDR,
			address: "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
			version: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			modules, err := c.AccountModules(tt.address, tt.version)
			if err != nil {
				t.Logf("get account modules for %s error: %s", tt.name, err.Error())
				return
			}

			b, _ := json.Marshal(modules)
			t.Logf("modules: %s", string(b))
		})
	}
}

func TestAptClient_AccountModuleById(t *testing.T) {
	tests := []struct {
		name     string
		rpc      string
		address  string
		moduleId string
		version  string
	}{
		{
			name:     "normal",
			rpc:      RPC_ADDR,
			address:  "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
			moduleId: "coin",
			version:  "7920450",
		},
		{
			name:     "address is null",
			rpc:      RPC_ADDR,
			address:  "",
			moduleId: "aptos_coin",
			version:  "7920450",
		},
		{
			name:     "moduleId is null",
			rpc:      RPC_ADDR,
			address:  "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
			moduleId: "",
			version:  "7920450",
		},
		{
			name:     "version is null",
			rpc:      RPC_ADDR,
			address:  "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
			moduleId: "aptos_coin",
			version:  "",
		},
		{
			name:     "params is null",
			rpc:      RPC_ADDR,
			address:  "",
			moduleId: "",
			version:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			module, err := c.AccountModuleById(tt.address, tt.moduleId, tt.version)
			if err != nil {
				t.Logf("account module by id for %s error: %s", tt.name, err)
				return
			}

			b, _ := json.Marshal(module)
			t.Logf("module: %s", string(b))
		})
	}
}

func TestAptClient_Transactions(t *testing.T) {
	tests := []struct {
		name  string
		rpc   string
		limit uint16
		start uint64
	}{
		{
			name:  "normal",
			rpc:   RPC_ADDR,
			limit: 25,
			start: 1,
		},
		{
			name:  "limit < 0",
			rpc:   RPC_ADDR,
			limit: 0,
			start: 1,
		},
		{
			name:  "start < 0",
			rpc:   RPC_ADDR,
			limit: 25,
			start: 1,
		},
		{
			name:  "limit < 0 && start < 0",
			rpc:   RPC_ADDR,
			limit: 0,
			start: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			transactions, err := c.Transactions(tt.limit, tt.start)
			if err != nil {
				t.Logf("get transaction for %s error: %s", tt.name, err.Error())
				return
			}

			b, _ := json.Marshal(transactions)
			t.Logf("transactions: %s", string(b))
		})
	}
}

func TestAptClient_TransactionsByAccount(t *testing.T) {
	tests := []struct {
		name    string
		rpc     string
		address string
		limit   uint16
		start   uint64
	}{
		{
			name:    "normal",
			rpc:     RPC_ADDR,
			address: "0x8b71b7d40de6ab3feea38c668bb3eba7152f6d45208b6d864c8587202e4d0c97",
			limit:   0,
			start:   0,
		},
		{
			name:    "address is null",
			rpc:     RPC_ADDR,
			address: "",
			limit:   0,
			start:   0,
		},
		{
			name:    "limit < 0",
			rpc:     RPC_ADDR,
			address: "0x8b71b7d40de6ab3feea38c668bb3eba7152f6d45208b6d864c8587202e4d0c97",
			limit:   0,
			start:   0,
		},
		{
			name:    "start < 0",
			rpc:     RPC_ADDR,
			address: "0x8b71b7d40de6ab3feea38c668bb3eba7152f6d45208b6d864c8587202e4d0c97",
			limit:   25,
			start:   1,
		},
		{
			name:    "limit < 0 && start < 0",
			rpc:     RPC_ADDR,
			address: "0x8b71b7d40de6ab3feea38c668bb3eba7152f6d45208b6d864c8587202e4d0c97",
			limit:   0,
			start:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			transactions, err := c.TransactionsByAccount(tt.address, tt.limit, tt.start)
			if err != nil {
				t.Logf("transaction by account for %s error: %s", tt.name, err.Error())
				return
			}
			b, _ := json.Marshal(transactions)
			t.Logf("transaction: %s", string(b))
		})
	}
}

func TestAptClient_TransactionByHash(t *testing.T) {
	tests := []struct {
		name string
		rpc  string
		hash string
	}{
		{
			name: "normal by hash",
			rpc:  MAINNET_RPC_ADDR,
			hash: "0x897ef35cea199183b30b84c5a6a690cd4f9e3f31ade4670dc12d50266230478a",
		},
		{
			name: "params null",
			rpc:  RPC_ADDR,
			hash: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			transaction, err := c.TransactionByHash(tt.hash)
			if err != nil {
				t.Logf("transaction by hashOrVersion %s error: %s", tt.name, err.Error())
				return
			}

			b, _ := json.Marshal(transaction)
			t.Logf("transaction: %s", string(b))
		})
	}
}

func TestAptClient_TransactionByVersion(t *testing.T) {
	tests := []struct {
		name    string
		rpc     string
		version uint64
	}{
		{
			name:    "normal by hash",
			rpc:     RPC_ADDR,
			version: 7935073,
		},
		{
			name:    "params null",
			rpc:     RPC_ADDR,
			version: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			transaction, err := c.TransactionByVersion(tt.version)
			if err != nil {
				t.Logf("transaction by hashOrVersion %s error: %s", tt.name, err.Error())
				return
			}

			b, _ := json.Marshal(transaction)
			t.Logf("transaction: %s", string(b))
		})
	}
}

func TestAptClient_SignMessage(t *testing.T) {
	tests := []struct {
		name    string
		rpc     string
		unSigTx *types.UnsignedTx
	}{
		{
			name: "normal",
			rpc:  RPC_ADDR,
			unSigTx: &types.UnsignedTx{
				Sender:          "0x593f8077f72f14e702f3b0fc0c362119b7c8c060282c3fb6e52311f525499f1a",
				SequenceNumber:  5,
				MaxGasAmount:    10,
				GasUnitPrice:    1,
				GasCurrencyCode: "",
				ExpirationTime:  uint64(time.Now().Unix()) + uint64(600),
				Payload: &types.EntryFunctionPayload{
					Type:          "entry_function_payload",
					Function:      "0x1::coin::transfer",
					TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
					Arguments:     []string{"0x6d829df49edf618de9002d16b03118f50cb0b22cb56901349720a07f6a5b10c5", strconv.FormatUint(100, 10)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			msg, err := c.SignMessage(tt.unSigTx)
			if err != nil {
				t.Logf("sign message for %s error: %s", tt.name, err.Error())
				return
			}

			b, _ := json.Marshal(msg)
			t.Logf("message: %s", string(b))
		})
	}
}

func TestAptClient_SendTransaction(t *testing.T) {
	tests := []struct {
		name        string
		rpc         string
		mnemonic    string
		index       int
		receiptAddr string
		amount      uint64
	}{
		{
			name:        "normal",
			rpc:         MAINNET_RPC_ADDR,
			mnemonic:    os.Getenv("KEY"),
			index:       1,
			receiptAddr: "0xf78ce06edced2ad1b471427485f945a22d9b5e21b50cbff4bdd25bcc90b716f9",
			amount:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			na := NewAptAccount(tt.mnemonic, "")
			account, err := na.AccountFromMnemonic(tt.index)
			if err != nil {
				t.Logf("new account from mnemonic error: %s", err)
				return
			}

			nonce, err := c.GetNonce(account.Address)
			if err != nil {
				t.Logf("get nonce for %s error: %s", tt.name, err.Error())
				return
			}

			unsignedTx, err := getTransferAptParams(account, tt.receiptAddr, tt.amount, nonce)
			if err != nil {
				t.Logf("getTransferAptParams for %s error: %s", tt.name, err.Error())
				return
			}

			signedTx, err := c.SignTransaction(account, unsignedTx)
			if err != nil {
				t.Logf("sign tx for %s error: %s", tt.name, err.Error())
				return
			}

			submitTx, err := c.SubmitTx(signedTx)
			if err != nil {
				t.Logf("submit transaction for %s error: %s", tt.name, err.Error())
				return
			}

			isSuccess, err := asyncTxStatus(submitTx.Hash)
			if err != nil {
				t.Logf("sync tx status for %s error: %s", tt.name, err.Error())
				return
			}

			t.Logf("send apt for %s from %s to %s result: %v, by hash: %s", tt.name, account.Address,
				tt.receiptAddr, isSuccess, "https://explorer.aptoslabs.com/txn/"+submitTx.Hash)
		})
	}
}

func TestAptClient_SimulateTx(t *testing.T) {
	tests := []struct {
		name        string
		rpc         string
		mnemonic    string
		index       int
		receiptAddr string
		amount      uint64
	}{
		{
			name:        "normal",
			rpc:         RPC_ADDR,
			mnemonic:    os.Getenv("KEY"),
			index:       0,
			receiptAddr: "0x6d829df49edf618de9002d16b03118f50cb0b22cb56901349720a07f6a5b10c5",
			amount:      1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewAptClient(tt.rpc)
			if err != nil {
				t.Logf("new apt client error: %s", err)
				return
			}

			na := NewAptAccount(tt.mnemonic, "")
			account, err := na.AccountFromMnemonic(tt.index)
			if err != nil {
				t.Logf("new account from mnemonic error: %s", err)
				return
			}

			nonce, err := c.GetNonce(account.Address)
			if err != nil {
				t.Logf("get balance for %s error: %s", tt.name, err.Error())
				return
			}

			unsignedTx, err := getTransferAptParams(account, tt.receiptAddr, tt.amount, nonce)
			if err != nil {
				t.Logf("getTransferAptParams for %s error: %s", tt.name, err.Error())
				return
			}

			signedTx, err := c.SignTransaction(account, unsignedTx)
			if err != nil {
				t.Logf("sign tx for %s error: %s", tt.name, err.Error())
				return
			}

			signedTx.Signature.Signature = "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
			tx, err := c.SimulateTx(signedTx)
			if err != nil {
				t.Logf("simulate transaction for %s error: %s", tt.name, err.Error())
				return
			}

			b, _ := json.Marshal(tx)
			t.Logf("simulate transaction result: %s", string(b))
		})
	}
}

func TestAptClient_EstimateGasPrice(t *testing.T) {
	c, err := NewAptClient(RPC_ADDR)
	if err != nil {
		t.Logf("new apt client error: %s", err)
		return
	}

	gasPrice, err := c.EstimateGasPrice()
	if err != nil {
		t.Logf("get estimate gas price error: %s", err.Error())
		return
	}
	t.Logf("gas price result: %d", gasPrice)
}

func TestAptClient_GetEventsByHandle(t *testing.T) {
	tests := []struct {
		name      string
		address   string
		handle    string
		fieldName string
		limit     uint16
		start     uint64
	}{
		{
			name:      "normal",
			address:   "0x8b71b7d40de6ab3feea38c668bb3eba7152f6d45208b6d864c8587202e4d0c97",
			handle:    "bbb",
			fieldName: "aaa",
			limit:     25,
			start:     1,
		},
	}

	c, err := NewAptClient(RPC_ADDR)
	if err != nil {
		t.Logf("new apt client error: %s", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, errEvent := c.GetEventsByHandle(tt.address, tt.handle, tt.fieldName, tt.limit, tt.start)
			if errEvent != nil {
				t.Logf("get estimate gas price error: %s", errEvent.Error())
				return
			}

			b, _ := json.Marshal(events)
			t.Logf("events result: %s", string(b))
		})
	}
}

func asyncTxStatus(txHash string) (bool, error) {
	count := 0
	c, _ := NewAptClient(MAINNET_RPC_ADDR)
	for {
		time.Sleep(1 * time.Second)
		txStatus, err := c.TransactionByHash(txHash)
		if err != nil {
			return false, err
		}

		if txStatus.Version != "0" {
			return txStatus.Success, nil
		}

		count++
		if count <= 10 {
			return false, fmt.Errorf("waiting for transaction %s timed out", txHash)
		}
	}

}

func genUnSignTx(account *types.AptAccount, nonce uint64, payload interface{}) (*types.UnsignedTx, error) {
	if account == nil {
		return nil, fmt.Errorf("account is null")
	}

	return &types.UnsignedTx{
		Sender:          account.Address,
		SequenceNumber:  nonce,
		MaxGasAmount:    2000,
		GasUnitPrice:    100,
		GasCurrencyCode: "XUS",
		ExpirationTime:  uint64(time.Now().Add(10 * time.Second).Unix()),
		Payload:         payload,
	}, nil
}

func getTransferAptParams(account *types.AptAccount, receiptAddr string, amount, nonce uint64) (
	*types.UnsignedTx,
	error,
) {
	payload := &types.EntryFunctionPayload{
		Type:          "entry_function_payload",
		Function:      "0x1::coin::transfer",
		TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
		Arguments:     []string{receiptAddr, strconv.FormatUint(amount, 10)},
	}

	unsignedTx, err := genUnSignTx(account, nonce, payload)
	if err != nil {
		return nil, err
	}
	return unsignedTx, nil
}
