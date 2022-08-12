## aptos-client

---
A Go library for Aptos chain SDK


### Requirements

- Go >= 1.18
- Aptos Node


### Install
```shell
go get github.com/threeandtwo/aptclient
```

### Iface
```go
AccountFromRandomKey() (*types.AptAccount, error)
AccountFromPrivateKey() (*types.AptAccount, error)
AccountFromMnemonic(index int) (*types.AptAccount, error)

Sign(msg []byte) []byte
Verify(sig, msg []byte) bool

LedgerInfo() (*types.LedgerInfo, error)
Account(address string) (*types.Account, error)
GetBalance(address string) (*big.Int, error)
GetNonce(address string) (uint64, error)
AccountResources(address, version string) ([]*types.AccountResource, error)
AccountResourceByType(address, resourceType, version string) (*types.AccountResource, error)
AccountModules(address, version string) ([]*types.AccountModule, error)
AccountModuleById(address, moduleID, version string) (*types.AccountModule, error)
Transactions(limit, start int) ([]*types.Transaction, error)
TransactionsByAccount(address string, limit, start int) ([]*types.Transaction, error)
Transaction(hashOrVersion string) (*types.Transaction, error)
SignMessage(unSigTx *types.UnsignedTx) (*types.SigningMessage, error)
SignTransaction(account *types.AptAccount, unsignedTx *types.UnsignedTx) (*types.SignedTx, error)
SubmitTx(signedTx *types.SignedTx) (*types.Transaction, error)
SimulateTx(signedTx *types.SignedTx) (*types.Transaction, error)
```

### Usage
| aptclient_test.go