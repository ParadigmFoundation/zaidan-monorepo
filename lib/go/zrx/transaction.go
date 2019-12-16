package zrx

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/eth"

	"github.com/0xProject/0x-mesh/ethereum"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core"
)

const ZeroExProtocolName = "0x Protocol"

const ZeroExProtocolVersion = "3.0.0"

// ZeroExTransaction represents 0x transaction (see ZEIP-18)
type ZeroExTransaction struct {
	Salt                  *big.Int
	ExpirationTimeSeconds *big.Int
	GasPrice              *big.Int
	SignerAddress         common.Address
	Data                  []byte

	hash *common.Hash
}

// MarshalJSON implements json.Marshaler
func (tx *ZeroExTransaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"salt":                  tx.Salt.String(),
		"expirationTimeSeconds": tx.ExpirationTimeSeconds.String(),
		"gasPrice":              tx.GasPrice.String(),
		"signerAddress":         strings.ToLower(tx.SignerAddress.Hex()),
		"data":                  hex.EncodeToString(tx.Data),
	})
}

// UnmarshalJSON implements json.Unmarshaler
func (tx *ZeroExTransaction) UnmarshalJSON(data []byte) error {
	var raw zeroExTransactionJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if err := tx.fromJSON(&raw); err != nil {
		return err
	}

	return nil
}

// ComputeHashForChainID calculates the 0x transaction hash for the provided chain ID.
// See https://github.com/0xProject/0x-protocol-specification/blob/master/v3/v3-specification.md#hashing-a-transaction
func (tx *ZeroExTransaction) ComputeHashForChainID(chainID int) (common.Hash, error) {
	if tx.hash != nil {
		return *tx.hash, nil
	}

	contractAddresses, err := ethereum.GetContractAddressesForChainID(chainID)
	if err != nil {
		return common.Hash{}, err
	}

	evmChainID := math.NewHexOrDecimal256(int64(chainID))
	domain := core.TypedDataDomain{
		Name:              ZeroExProtocolName,
		Version:           ZeroExProtocolVersion,
		ChainId:           evmChainID,
		VerifyingContract: contractAddresses.Exchange.Hex(),
	}

	typedData := core.TypedData{
		Types:       EIP712Types,
		PrimaryType: TypeZeroExTransaction,
		Domain:      domain,
		Message:     tx.Map(),
	}

	domainSeparator, err := typedData.HashStruct(TypeEIP712Domain, typedData.Domain.Map())
	if err != nil {
		return common.Hash{}, err
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return common.Hash{}, err
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hashBytes := eth.Keccak256(rawData)
	hash := common.BytesToHash(hashBytes)
	tx.hash = &hash
	return hash, nil
}

// Map returns the transaction as an un-typed map (useful when hashing)
func (tx *ZeroExTransaction) Map() map[string]interface{} {
	return map[string]interface{}{
		"salt":                  tx.Salt.String(),
		"expirationTimeSeconds": tx.ExpirationTimeSeconds.String(),
		"gasPrice":              tx.GasPrice.String(),
		"signerAddress":         tx.SignerAddress.Hex(),
		"data":                  tx.Data,
	}
}

// ResetHash returns the cached transaction hash to nil
func (tx *ZeroExTransaction) ResetHash() {
	tx.hash = nil
}

// set a 0x transaction values from a JSON representation
func (tx *ZeroExTransaction) fromJSON(ztx *zeroExTransactionJSON) error {
	salt, ok := new(big.Int).SetString(ztx.Salt, 10)
	if !ok {
		return errors.New(`unable to unmarshal value for "salt"`)
	}

	expirationTimeSeconds, ok := new(big.Int).SetString(ztx.ExpirationTimeSeconds, 10)
	if !ok {
		return errors.New(`unable to unmarshal value for "expirationTimeSeconds"`)
	}

	gasPrice, ok := new(big.Int).SetString(ztx.GasPrice, 10)
	if !ok {
		return errors.New(`unable to unmarshal value for "gasPrice"`)
	}

	tx.Salt = salt
	tx.ExpirationTimeSeconds = expirationTimeSeconds
	tx.GasPrice = gasPrice
	tx.SignerAddress = common.HexToAddress(ztx.SignerAddress)

	if ztx.Data[:2] == "0x" {
		tx.Data = common.Hex2Bytes(ztx.Data[2:])
	} else {
		tx.Data = common.Hex2Bytes(ztx.Data)
	}

	return nil
}

// used to assist in un-marshalling 0x transactions
type zeroExTransactionJSON struct {
	Salt                  string `json:"salt"`
	ExpirationTimeSeconds string `json:"expirationTimeSeconds"`
	GasPrice              string `json:"gasPrice"`
	SignerAddress         string `json:"signerAddress"`
	Data                  string `json:"data"`
}
