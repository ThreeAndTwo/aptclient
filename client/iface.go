package client

import (
	"github.com/threeandtwo/aptclient/types"
	"math/big"
)

type (
	IAccount interface {
		AccountFromRandomKey() (*types.AptAccount, error)
		AccountFromPrivateKey() (*types.AptAccount, error)
		AccountFromMnemonic(index int) (*types.AptAccount, error)
		GetAptAccount(index int) (*types.AptAccount, error)

		Sign(msg []byte) []byte
		Verify(sig, msg []byte) bool
	}

	IClient interface {
		NodeHealth(durationSecs uint32) (string, error)
		LedgerInfo() (*types.LedgerInfo, error)

		BlockByHeight(blockHeight uint64, withTxs types.BlockWithTxs) (*types.Block, error)
		BlockByVersion(version uint64, withTxs types.BlockWithTxs) (*types.Block, error)

		Account(address string) (*types.Account, error)
		GetBalance(address string) (*big.Int, error)
		GetNonce(address string) (uint64, error)
		AccountResources(address, version string) ([]*types.AccountResource, error)
		AccountResourceByType(address, resourceType, version string) (*types.AccountResource, error)
		AccountModules(address, version string) ([]*types.AccountModule, error)
		AccountModuleById(address, moduleID, version string) (*types.AccountModule, error)

		Transactions(limit, start int) ([]*types.Transaction, error)
		TransactionsByAccount(address string, limit, start int) ([]*types.Transaction, error)
		TransactionByHash(hash string) (*types.Transaction, error)
		TransactionByVersion(version uint64) (*types.Transaction, error)
		SignMessage(unSigTx *types.UnsignedTx) (*types.SigningMessage, error)
		SignTransaction(account *types.AptAccount, unsignedTx *types.UnsignedTx) (*types.SignedTx, error)
		SubmitTx(signedTx *types.SignedTx) (*types.Transaction, error)
		SimulateTx(signedTx *types.SignedTx) ([]*types.SimulateTx, error)
	}
)
