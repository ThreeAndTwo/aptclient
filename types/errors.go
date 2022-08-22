package types

import "errors"

var (
	ErrNotPrivateKeyTy = errors.New("params not privateKey")
	ErrNotNoneTy       = errors.New("params mismatched")
	ErrNotMnemonicTy   = errors.New("params not mnemonic")
	ErrMnemonicIndex   = errors.New("index must be ge 0 for mnemonic")
	ErrMnemonicCount   = errors.New("mnemonic count should be 15 | 18 | 21 | 24")

	ErrRpcNull          = errors.New("rps address is null")
	ErrAddressNull      = errors.New("address is null")
	ErrAddressLen       = errors.New("address length != 66")
	ErrResourceTypeNull = errors.New("resource type is null")
	ErrParsedValue      = errors.New("value type is mismatched")
	ErrModuleIdNull     = errors.New("moduleId is null")
	ErrHashNull         = errors.New("hash is null")
	ErrSignNull         = errors.New("signature is null")
	ErrPayloadNull      = errors.New("payload is null")
	ErrRequestRpc       = errors.New("request REST API error")
)
