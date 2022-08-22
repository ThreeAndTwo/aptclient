package types

import "crypto/ed25519"

type KeyTy int

const (
	MnemonicTy KeyTy = iota
	PrivateTy
	NoneTy
)

type BlockWithTxs string

const (
	FalseTy BlockWithTxs = "false"
	TrueTy               = "true"
)

const (
	AptResourceTy = "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>"
	AptAccountTy  = "0x1::account::Account"
	Ed25519       = "ed25519_signature"
)

type NodeHealth struct {
	Message string `json:"message"`
}

type LedgerInfo struct {
	ChainID             int32  `json:"chain_id"`
	Epoch               string `json:"epoch"`
	LedgerVersion       uint   `json:"ledger_version,string"`
	OldestLedgerVersion string `json:"oldest_ledger_version"`
	LedgerTimestamp     uint64 `json:"ledger_timestamp,string"`
	NodeRole            string `json:"node_role"`
}

type Block struct {
	BlockHeight    string               `json:"block_height"`
	BlockHash      string               `json:"block_hash"`
	BlockTimestamp string               `json:"block_timestamp"`
	FirstVersion   string               `json:"first_version"`
	LastVersion    string               `json:"last_version"`
	Transactions   []TransactionInBlock `json:"transactions"`
}

type TransactionInBlock struct {
	Type                    string                `json:"type"`
	Hash                    string                `json:"hash"`
	Sender                  string                `json:"sender"`
	SequenceNumber          string                `json:"sequence_number"`
	MaxGasAmount            string                `json:"max_gas_amount"`
	GasUnitPrice            string                `json:"gas_unit_price"`
	ExpirationTimestampSecs string                `json:"expiration_timestamp_secs"`
	Payload                 *EntryFunctionPayload `json:"payload"`
	Signature               *TxSignature          `json:"signature"`
}

type AptAccount struct {
	Address    string
	PublicKey  string
	PrivateKey ed25519.PrivateKey
	AuthKey    string
}

type Account struct {
	SequenceNumber uint64 `json:"sequence_number,string"`
	AuthKey        string `json:"authentication_key"`
}

type AccountResource struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type AccountModule struct {
	ByteCode string         `json:"bytecode"`
	ABI      *MoveModuleABI `json:"abi"`
}

type MoveModuleABI struct {
	Address string   `json:"address"`
	Name    string   `json:"string"`
	Friends []string `json:"friends"`
}

type MoveFunction struct {
	Name              string        `json:"name"`
	Visibility        string        `json:"visibility"`
	GenericTypeParams []interface{} `json:"generic_type_params"`
	Params            []string      `json:"params"`
	Returns           []string      `json:"returns"`
}

type MoveStruct struct {
	Name              string        `json:"name"`
	IsNative          bool          `json:"is_native"`
	Abilities         []string      `json:"abilities"`
	GenericTypeParams []interface{} `json:"generic_type_params"`
}

type MoveStructField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Transaction struct {
	Type                string           `json:"type"`
	Events              []TxEvents       `json:"events"`
	Payload             *WriteSetPayload `json:"payload"`
	Version             string           `json:"version"`
	Hash                string           `json:"hash"`
	StateRootHash       string           `json:"state_root_hash"`
	EventRootHash       string           `json:"event_root_hash"`
	GasUsed             uint64           `json:"gas_used,string"`
	Success             bool             `json:"success"`
	VMStatus            string           `json:"vm_status"`
	AccumulatorRootHash string           `json:"accumulator_root_hash"`
	Changes             []TxChange       `json:"changes"`
}

type SimulateTx struct {
	Type                    string               `json:"type"`
	Version                 string               `json:"version"`
	Hash                    string               `json:"hash"`
	StateRootHash           string               `json:"state_root_hash"`
	EventRootHash           string               `json:"event_root_hash"`
	GasUsed                 string               `json:"gas_used"`
	Success                 bool                 `json:"success"`
	VMStatus                string               `json:"vm_status"`
	AccumulatorRootHash     string               `json:"accumulator_root_hash"`
	Changes                 []TxChange           `json:"changes"`
	Sender                  string               `json:"sender"`
	SequenceNumber          string               `json:"sequence_number"`
	MaxGasAmount            string               `json:"max_gas_amount"`
	GasUnitPrice            string               `json:"gas_unit_price"`
	ExpirationTimestampSecs string               `json:"expiration_timestamp_secs"`
	Payload                 EntryFunctionPayload `json:"payload"`
	Signature               TxSignature          `json:"signature"`
	Events                  []TxEvents           `json:"events"`
	Timestamp               string               `json:"timestamp"`
}

type TxChange struct {
	Type         string `json:"type"`
	StateKeyHash string `json:"state_key_hash"`
	Address      string `json:"address"`
	Module       string `json:"module"`
}

type TxEvents struct {
	Key            string      `json:"key"`
	SequenceNumber string      `json:"sequence_number"`
	Type           string      `json:"type"`
	Data           TxEventData `json:"data"`
}

type TxEventData struct {
	Created string `json:"created"`
	RoleId  string `json:"role_id"`
}

type WriteSetPayload struct {
	Type     string       `json:"type"`
	WriteSet WriteSetData `json:"write_set"`
}

type WriteSetData struct {
	Type      string      `json:"type"`
	ExecuteAs string      `json:"execute_as"`
	Script    interface{} `json:"script"`
}

type UnsignedTx struct {
	Sender          string      `json:"sender,string"`
	SequenceNumber  uint64      `json:"sequence_number,string"`
	MaxGasAmount    uint64      `json:"max_gas_amount,string"`
	GasUnitPrice    uint64      `json:"gas_unit_price,string"`
	GasCurrencyCode string      `json:"gas_currency_code"`
	ExpirationTime  uint64      `json:"expiration_timestamp_secs,string"`
	Payload         interface{} `json:"payload"`
}

type SignedTx struct {
	*UnsignedTx
	Signature *TxSignature `json:"signature"`
}

type TxSignature struct {
	Type      string `json:"type"`
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

type EntryFunctionPayload struct {
	Type          string   `json:"type"`
	Function      string   `json:"function"`
	TypeArguments []string `json:"type_arguments"`
	Arguments     []string `json:"arguments"`
}

type SigningMessage struct {
	Message string `json:"message"`
}

type ExceptionMsg struct {
	Message       string `json:"message"`
	Code          string `json:"error_code"`
	LedgerVersion string `json:"ledger_version"`
}
