package client

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/threeandtwo/aptclient/types"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

type AptClient struct {
	rpc string
}

func (a *AptClient) NodeHealth(durationSecs uint32) (string, error) {
	rpc := fmt.Sprintf("%s/-/healthy?duration_secs=%d", a.rpc, durationSecs)
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return "", err
	}
	return req, err
}

func (a *AptClient) LedgerInfo() (*types.LedgerInfo, error) {
	rpc := fmt.Sprintf("%s/", a.rpc)

	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	_info := &types.LedgerInfo{}
	err = json.Unmarshal([]byte(req), _info)
	return _info, err
}

func (a *AptClient) BlockByHeight(blockHeight uint64, withTxs types.BlockWithTxs) (*types.Block, error) {
	if withTxs == "" {
		withTxs = types.FalseTy
	}

	rpc := fmt.Sprintf("%s/blocks/by_height/%d?with_transactions=%s", a.rpc, blockHeight, withTxs)
	return a.getBlockInfo(rpc)
}

func (a *AptClient) BlockByVersion(version uint64, withTxs types.BlockWithTxs) (*types.Block, error) {
	if withTxs == "" {
		withTxs = types.FalseTy
	}

	rpc := fmt.Sprintf("%s/blocks/by_version/%d?with_transactions=%s", a.rpc, version, withTxs)
	return a.getBlockInfo(rpc)
}

func (a *AptClient) getBlockInfo(rpc string) (*types.Block, error) {
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	_block := &types.Block{}
	err = json.Unmarshal([]byte(req), _block)
	return _block, err
}

func (a *AptClient) Account(address string) (*types.Account, error) {
	isCheck, err := checkAccount(address)
	if isCheck {
		return nil, err
	}

	rpc := fmt.Sprintf("%s/accounts/%s", a.rpc, address)
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	_account := &types.Account{}
	err = json.Unmarshal([]byte(req), _account)
	return _account, err
}

func checkAccount(address string) (bool, error) {
	if address == "" {
		return true, types.ErrAddressNull
	}

	if len(address) != 66 {
		return true, types.ErrAddressLen
	}
	return false, nil
}

func (a *AptClient) GetBalance(address string) (*big.Int, error) {
	isCheck, err := checkAccount(address)
	if isCheck {
		return nil, err
	}

	res, err := a.AccountResourceByType(address, types.AptResourceTy, "")
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, types.ErrResourceTypeNull
	}

	if _, ok := res.Data["coin"]; !ok || reflect.TypeOf(res.Data["coin"]).Kind() != reflect.Map {
		return nil, types.ErrParsedValue
	}

	v := res.Data["coin"].(map[string]interface{})

	_v, err := strconv.ParseInt(v["value"].(string), 10, 0)
	if err != nil {
		return nil, err
	}
	return big.NewInt(_v), nil
}

func (a *AptClient) GetNonce(address string) (uint64, error) {
	isCheck, err := checkAccount(address)
	if isCheck {
		return 0, err
	}

	res, err := a.AccountResourceByType(address, types.AptAccountTy, "")
	if err != nil {
		return 0, err
	}

	if res.Data == nil {
		return 0, types.ErrResourceTypeNull
	}

	if _, ok := res.Data["sequence_number"]; !ok || reflect.TypeOf(res.Data["sequence_number"]).Kind() != reflect.String {
		return 0, types.ErrParsedValue
	}

	nonce := res.Data["sequence_number"].(string)
	return strconv.ParseUint(nonce, 10, 64)
}

func (a *AptClient) AccountResources(address, version string) ([]*types.AccountResource, error) {
	isCheck, err := checkAccount(address)
	if isCheck {
		return nil, err
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

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	err = json.Unmarshal([]byte(req), &_as)
	return _as, err
}

func (a *AptClient) AccountResourceByType(address, resourceType, version string) (*types.AccountResource, error) {
	isCheck, err := checkAccount(address)
	if isCheck {
		return nil, err
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

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	err = json.Unmarshal([]byte(req), &_as)
	return _as, err
}

func (a *AptClient) AccountModules(address, version string) ([]*types.AccountModule, error) {
	isCheck, err := checkAccount(address)
	if isCheck {
		return nil, err
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

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	err = json.Unmarshal([]byte(req), &_am)
	return _am, err
}

func (a *AptClient) AccountModuleById(address, moduleName, version string) (*types.AccountModule, error) {
	isCheck, err := checkAccount(address)
	if isCheck {
		return nil, err
	}

	if moduleName == "" {
		return nil, types.ErrModuleIdNull
	}

	rpc := ""
	if version == "" {
		rpc = fmt.Sprintf("%s/accounts/%s/module/%s", a.rpc, address, moduleName)
	} else {
		rpc = fmt.Sprintf("%s/accounts/%s/module/%s?version=%s", a.rpc, address, moduleName, version)
	}

	var _am *types.AccountModule
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	err = json.Unmarshal([]byte(req), &_am)
	return _am, err
}

func (a *AptClient) Transactions(limit uint16, start uint64) ([]*types.Transaction, error) {
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

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	err = json.Unmarshal([]byte(req), &txs)
	return txs, err
}

func (a *AptClient) TransactionsByAccount(address string, limit uint16, start uint64) ([]*types.Transaction, error) {
	if limit <= 0 {
		limit = 25
	}

	if start <= 0 {
		start = 1
	}

	isCheck, err := checkAccount(address)
	if isCheck {
		return nil, err
	}

	rpc := fmt.Sprintf("%s/accounts/%s/transactions?limit=%d&start=%d", a.rpc, address, limit, start)

	var txs []*types.Transaction
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	err = json.Unmarshal([]byte(req), &txs)
	return txs, err
}

func (a *AptClient) TransactionByHash(hash string) (*types.Transaction, error) {
	if hash == "" {
		return nil, types.ErrHashNull
	}

	rpc := fmt.Sprintf("%s/transactions/by_hash/%s", a.rpc, hash)
	return a.transaction(rpc)
}

func (a *AptClient) TransactionByVersion(version uint64) (*types.Transaction, error) {
	rpc := fmt.Sprintf("%s/transactions/by_version/%d", a.rpc, version)
	return a.transaction(rpc)
}

func (a *AptClient) transaction(rpc string) (*types.Transaction, error) {
	var tx *types.Transaction
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	err = json.Unmarshal([]byte(req), &tx)
	return tx, err
}

func (a *AptClient) SignMessage(unSigTx *types.UnsignedTx) (*types.SigningMessage, error) {
	if unSigTx.Payload == nil {
		return nil, types.ErrPayloadNull
	}

	rpc := fmt.Sprintf("%s/transactions/encode_submission", a.rpc)
	unsignedMap := initUnSigMap(unSigTx)

	sigMsg := &types.SigningMessage{}
	req, err := a.connClient(rpc, unsignedMap).Request(PostTy)
	if err != nil {
		return nil, err
	}

	if req == "" {
		return nil, types.ErrSignNull
	}

	sigMsg.Message = strings.Trim(req, `"`)
	if sigMsg.Message[0:2] != "0x" {
		return nil, fmt.Errorf(req)
	}

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

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	err = json.Unmarshal([]byte(req), &tx)
	return tx, err
}

func (a *AptClient) SimulateTx(signedTx *types.SignedTx) ([]*types.SimulateTx, error) {
	rpc := fmt.Sprintf("%s/transactions/simulate", a.rpc)
	signedMap := initSigTx(signedTx)

	var tx []*types.SimulateTx
	req, err := a.connClient(rpc, signedMap).Request(PostTy)
	if err != nil {
		return nil, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	err = json.Unmarshal([]byte(req), &tx)
	return tx, err
}

func (a *AptClient) SubmitBatchTx(signedTxs []*types.SignedTx) error {
	rpc := fmt.Sprintf("%s/transactions/batch", a.rpc)

	var batchedSignedTx map[string]interface{}
	for k, signedTx := range signedTxs {
		signedMap := initSigTx(signedTx)
		batchedSignedTx[strconv.Itoa(k)] = signedMap
	}

	var tx []*types.SimulateTx
	req, err := a.connClient(rpc, batchedSignedTx).Request(PostTy)
	if err != nil {
		return err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return fmt.Errorf(errDesc)
	}

	err = json.Unmarshal([]byte(req), &tx)
	return err
}

func (a *AptClient) EstimateGasPrice() (uint64, error) {
	rpc := fmt.Sprintf("%s/estimate_gas_price", a.rpc)
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return 0, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return 0, fmt.Errorf(errDesc)
	}

	var gp types.EstimateGasPrice
	err = json.Unmarshal([]byte(req), &gp)
	return gp.GasEstimate, err
}

// Deprecated
func (a *AptClient) GetEventsByKey(key string, limit uint16, start uint64) ([]*types.Event, error) {
	if key == "" {
		return nil, fmt.Errorf("key should be null")
	}

	rpc := fmt.Sprintf("%s/events/%s?limit=%d&start=%d", a.rpc, key, limit, start)

	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	var events []*types.Event
	err = json.Unmarshal([]byte(req), &events)
	return events, err
}

func (a *AptClient) GetEventsByCreationNumber(address, creationNumber string, limit, start uint64) ([]*types.Event, error) {
	if address == "" || creationNumber == "" {
		return nil, fmt.Errorf("address | handle | fieldName is null, plz check it")
	}

	rpc := fmt.Sprintf("%s/accounts/%s/events/%s?limit=%d&start=%d", a.rpc, address, creationNumber, limit, start)
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	var events []*types.Event
	err = json.Unmarshal([]byte(req), &events)
	return events, err

}

func (a *AptClient) GetEventsByHandle(address, handle, fieldName string, limit uint16, start uint64) ([]*types.Event, error) {
	if address == "" || handle == "" || fieldName == "" {
		return nil, fmt.Errorf("address | handle | fieldName is null, plz check it")
	}

	rpc := fmt.Sprintf("%s/accounts/%s/events/%s/%s?limit=%d&start=%d", a.rpc, address, handle, fieldName, limit, start)
	req, err := a.connClient(rpc, nil).Request(GetTy)
	if err != nil {
		return nil, err
	}

	hasE, errDesc := hasExceptionForResp(req)
	if hasE {
		return nil, fmt.Errorf(errDesc)
	}

	var events []*types.Event
	err = json.Unmarshal([]byte(req), &events)
	return events, err
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

func hasExceptionForResp(msg string) (bool, string) {
	exMsg := &types.ExceptionMsg{}
	errMsg := ""

	if json.Unmarshal([]byte(msg), exMsg) != nil {
		return false, ""
	}

	if exMsg.Message == "" {
		return false, ""
	}

	errMsg = types.ErrRequestRpc.Error() + ": " + exMsg.Message
	return true, errMsg
}

func initSigTx(signedTx *types.SignedTx) map[string]interface{} {
	signedMap := make(map[string]interface{})
	signedMap["sender"] = signedTx.Sender
	signedMap["sequence_number"] = fmt.Sprintf("%d", signedTx.SequenceNumber)
	signedMap["max_gas_amount"] = fmt.Sprintf("%d", signedTx.MaxGasAmount)
	signedMap["gas_unit_price"] = fmt.Sprintf("%d", signedTx.GasUnitPrice)
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
	unsignedMap["expiration_timestamp_secs"] = fmt.Sprintf("%d", unSigTx.ExpirationTime)
	unsignedMap["payload"] = unSigTx.Payload
	return unsignedMap
}
