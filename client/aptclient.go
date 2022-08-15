package client

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/threeandtwo/aptclient/types"
	"math/big"
	"strconv"
	"strings"
)

type AptClient struct {
	rpc string
}

func (a *AptClient) LedgerInfo() (*types.LedgerInfo, error) {
	rpc := fmt.Sprintf("%s/", a.rpc)

	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	_info := &types.LedgerInfo{}
	err = json.Unmarshal([]byte(req), _info)
	return _info, err
}

func (a *AptClient) Account(address string) (*types.Account, error) {
	if checkAccount(address) {
		return nil, types.ErrAddressNull
	}

	rpc := fmt.Sprintf("%s/accounts/%s", a.rpc, address)
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	_account := &types.Account{}
	err = json.Unmarshal([]byte(req), _account)
	return _account, err
}

func checkAccount(address string) bool {
	return address == "" || len(address) != 66
}

func (a *AptClient) GetBalance(address string) (*big.Int, error) {
	if checkAccount(address) {
		return nil, types.ErrAddressNull
	}

	res, err := a.AccountResourceByType(address, types.AptResourceTy, "")
	if err != nil {
		return nil, err
	}
	v := res.Data["coin"].(map[string]interface{})

	_v, err := strconv.ParseInt(v["value"].(string), 10, 0)
	if err != nil {
		return nil, err
	}
	return big.NewInt(_v), nil
}

func (a *AptClient) GetNonce(address string) (uint64, error) {
	if checkAccount(address) {
		return 0, types.ErrAddressNull
	}

	res, err := a.AccountResourceByType(address, types.AptAccountTy, "")
	if err != nil {
		return 0, err
	}
	nonce := res.Data["sequence_number"].(string)
	return strconv.ParseUint(nonce, 10, 64)
}

func (a *AptClient) AccountResources(address, version string) ([]*types.AccountResource, error) {
	if checkAccount(address) {
		return nil, types.ErrAddressNull
	}

	rpc := ""
	if version == "" {
		rpc = fmt.Sprintf("%s/accounts/%s/resources", a.rpc, address)
	} else {
		rpc = fmt.Sprintf("%s/accounts/%s/resources?version=%s", a.rpc, address, version)
	}

	var _as []*types.AccountResource
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	err = json.Unmarshal([]byte(req), &_as)
	return _as, err
}

func (a *AptClient) AccountResourceByType(address, resourceType, version string) (*types.AccountResource, error) {
	if checkAccount(address) {
		return nil, types.ErrAddressNull
	}

	if resourceType == "" {
		return nil, types.ErrResourceTypeNull
	}

	rpc := ""
	if version == "" {
		rpc = fmt.Sprintf("%s/accounts/%s/resource/%s", a.rpc, address, resourceType)
	} else {
		rpc = fmt.Sprintf("%s/accounts/%s/resource/%s?version=%s", a.rpc, address, resourceType, version)
	}

	var _as *types.AccountResource
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	err = json.Unmarshal([]byte(req), &_as)
	return _as, err
}

func (a *AptClient) AccountModules(address, version string) ([]*types.AccountModule, error) {
	if checkAccount(address) {
		return nil, types.ErrAddressNull
	}

	rpc := ""
	if version == "" {
		rpc = fmt.Sprintf("%s/accounts/%s/modules", a.rpc, address)
	} else {
		rpc = fmt.Sprintf("%s/accounts/%s/modules?version=%s", a.rpc, address, version)
	}

	var _am []*types.AccountModule
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	err = json.Unmarshal([]byte(req), &_am)
	return _am, err
}

func (a *AptClient) AccountModuleById(address, moduleId, version string) (*types.AccountModule, error) {
	if checkAccount(address) {
		return nil, types.ErrAddressNull
	}

	if moduleId == "" {
		return nil, types.ErrModuleIdNull
	}

	rpc := ""
	if version == "" {
		rpc = fmt.Sprintf("%s/accounts/%s/module/%s", a.rpc, address, moduleId)
	} else {
		rpc = fmt.Sprintf("%s/accounts/%s/module/%s?version=%s", a.rpc, address, moduleId, version)
	}

	var _am *types.AccountModule
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	err = json.Unmarshal([]byte(req), &_am)
	return _am, err
}

func (a *AptClient) Transactions(limit, start int) ([]*types.Transaction, error) {
	if limit <= 0 {
		limit = 25
	}

	if start <= 0 {
		start = 1
	}

	rpc := fmt.Sprintf("%s/transactions?limit=%d&start=%d", a.rpc, limit, start)

	var txs []*types.Transaction
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	err = json.Unmarshal([]byte(req), &txs)
	return txs, err
}

func (a *AptClient) TransactionsByAccount(address string, limit, start int) ([]*types.Transaction, error) {
	if limit <= 0 {
		limit = 25
	}

	if start <= 0 {
		start = 1
	}

	if checkAccount(address) {
		return nil, types.ErrAddressNull
	}

	rpc := fmt.Sprintf("%s/accounts/%s/transactions?limit=%d&start=%d", a.rpc, address, limit, start)

	var txs []*types.Transaction
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	err = json.Unmarshal([]byte(req), &txs)
	return txs, err
}

func (a *AptClient) Transaction(hashOrVersion string) (*types.Transaction, error) {
	if hashOrVersion == "" {
		return nil, types.ErrHashNull
	}
	rpc := fmt.Sprintf("%s/transactions/%s", a.rpc, hashOrVersion)

	var tx *types.Transaction
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	err = json.Unmarshal([]byte(req), &tx)
	return tx, err
}

func (a *AptClient) SignMessage(unSigTx *types.UnsignedTx) (*types.SigningMessage, error) {
	rpc := fmt.Sprintf("%s/transactions/signing_message", a.rpc)
	unsignedMap := initUnSigMap(unSigTx)

	var sigMsg *types.SigningMessage
	req, err := a.connClient(rpc, unsignedMap).Request(PostTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	err = json.Unmarshal([]byte(req), &sigMsg)
	return sigMsg, err
}

func (a *AptClient) SignTransaction(account *types.AptAccount, unsignedTx *types.UnsignedTx) (*types.SignedTx, error) {
	msg, err := a.SignMessage(unsignedTx)
	if err != nil {
		return nil, err
	}

	hexMsg, err := hex.DecodeString(msg.Message[2:])
	if err != nil {
		return nil, err
	}

	sig := &types.TxSignature{
		Type:      types.Ed25519,
		PublicKey: account.PublicKey,
		Signature: hex.EncodeToString(ed25519.Sign(account.PrivateKey, hexMsg)),
	}

	return &types.SignedTx{
		UnsignedTx: unsignedTx,
		Signature:  sig,
	}, nil
}

func (a *AptClient) SubmitTx(signedTx *types.SignedTx) (*types.Transaction, error) {
	rpc := fmt.Sprintf("%s/transactions", a.rpc)
	signedMap := initSigTx(signedTx)

	var tx *types.Transaction
	req, err := a.connClient(rpc, signedMap).Request(PostTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	err = json.Unmarshal([]byte(req), &tx)
	return tx, err
}

func (a *AptClient) SimulateTx(signedTx *types.SignedTx) (*types.Transaction, error) {
	rpc := fmt.Sprintf("%s/transactions/simulate", a.rpc)
	signedMap := initSigTx(signedTx)

	var tx *types.Transaction
	req, err := a.connClient(rpc, signedMap).Request(PostTy)
	if err != nil {
		return nil, err
	}

	if hasExceptionForResp(req) {
		return nil, fmt.Errorf(req)
	}

	err = json.Unmarshal([]byte(req), &tx)
	return tx, err
}

func NewAptClient(rpc string) (*AptClient, error) {
	if rpc == "" {
		return nil, types.ErrRpcNull
	}

	rpc = fmtRpc(rpc)
	return &AptClient{rpc: rpc}, nil
}

func fmtRpc(rpc string) string {
	if rpc[len(rpc)-1:] == "/" {
		return strings.Trim(rpc[:len(rpc)-1], " ")
	}
	return strings.Trim(rpc, " ")
}

func initHeader() map[string]string {
	header := make(map[string]string)
	header["content-type"] = "application/json"
	return header
}

func (a *AptClient) connClient(url string, params map[string]interface{}) *Net {
	return NewNet(url, initHeader(), params)
}

func hasExceptionForResp(msg string) bool {
	exMsg := &types.ExceptionMsg{}

	if json.Unmarshal([]byte(msg), exMsg) != nil {
		return false
	}

	if exMsg.Code == 0 {
		return false
	}
	return true
}

func initSigTx(signedTx *types.SignedTx) map[string]interface{} {
	signedMap := make(map[string]interface{})
	signedMap["sender"] = signedTx.Sender
	signedMap["sequence_number"] = fmt.Sprintf("%d", signedTx.SequenceNumber)
	signedMap["max_gas_amount"] = fmt.Sprintf("%d", signedTx.MaxGasAmount)
	signedMap["gas_unit_price"] = fmt.Sprintf("%d", signedTx.GasUnitPrice)
	signedMap["gas_currency_code"] = signedTx.GasCurrencyCode
	signedMap["expiration_timestamp_secs"] = fmt.Sprintf("%d", signedTx.ExpirationTime)
	signedMap["payload"] = signedTx.Payload
	signedMap["signature"] = signedTx.Signature
	return signedMap
}

func initUnSigMap(unSigTx *types.UnsignedTx) map[string]interface{} {
	unsignedMap := make(map[string]interface{})
	unsignedMap["sender"] = unSigTx.Sender
	unsignedMap["sequence_number"] = fmt.Sprintf("%d", unSigTx.SequenceNumber)
	unsignedMap["max_gas_amount"] = fmt.Sprintf("%d", unSigTx.MaxGasAmount)
	unsignedMap["gas_unit_price"] = fmt.Sprintf("%d", unSigTx.GasUnitPrice)
	unsignedMap["gas_currency_code"] = unSigTx.GasCurrencyCode
	unsignedMap["expiration_timestamp_secs"] = fmt.Sprintf("%d", unSigTx.ExpirationTime)
	unsignedMap["payload"] = unSigTx.Payload
	return unsignedMap
}
