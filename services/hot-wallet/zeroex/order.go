package zeroex

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/math"

	"github.com/ethereum/go-ethereum/common"

	gethsigner "github.com/ethereum/go-ethereum/signer/core"
)

const addressHexLength = 42

// SignatureType represents the type of 0x signature encountered
type SignatureType uint8

// SignatureType values
const (
	IllegalSignature SignatureType = iota
	InvalidSignature
	EIP712Signature
	EthSignSignature
	WalletSignature
	ValidatorSignature
	PreSignedSignature
	EIP1271WalletSignature
	NSignatureTypesSignature
)

// Order represents an unsigned 0x order
type Order struct {
	ChainID               *big.Int       `json:"chainId"`
	ExchangeAddress       common.Address `json:"exchangeAddress"`
	MakerAddress          common.Address `json:"makerAddress"`
	MakerAssetData        []byte         `json:"makerAssetData"`
	MakerFeeAssetData     []byte         `json:"makerFeeAssetData"`
	MakerAssetAmount      *big.Int       `json:"makerAssetAmount"`
	MakerFee              *big.Int       `json:"makerFee"`
	TakerAddress          common.Address `json:"takerAddress"`
	TakerAssetData        []byte         `json:"takerAssetData"`
	TakerFeeAssetData     []byte         `json:"takerFeeAssetData"`
	TakerAssetAmount      *big.Int       `json:"takerAssetAmount"`
	TakerFee              *big.Int       `json:"takerFee"`
	SenderAddress         common.Address `json:"senderAddress"`
	FeeRecipientAddress   common.Address `json:"feeRecipientAddress"`
	ExpirationTimeSeconds *big.Int       `json:"expirationTimeSeconds"`
	Salt                  *big.Int       `json:"salt"`

	// Cache hash for performance
	hash *common.Hash
}

// SignedOrder represents a signed 0x order
type SignedOrder struct {
	Order
	Signature []byte `json:"signature"`
}

// ResetHash resets the cached order hash. Usually only required for testing.
func (o *Order) ResetHash() {
	o.hash = nil
}

// ComputeOrderHash computes a 0x order hash
func (o *Order) ComputeOrderHash() (common.Hash, error) {
	if o.hash != nil {
		return *o.hash, nil
	}

	chainID := math.NewHexOrDecimal256(o.ChainID.Int64())
	var domain = gethsigner.TypedDataDomain{
		Name:              "0x Protocol",
		Version:           "3.0.0",
		ChainId:           chainID,
		VerifyingContract: o.ExchangeAddress.Hex(),
	}

	var message = map[string]interface{}{
		"makerAddress":          o.MakerAddress.Hex(),
		"takerAddress":          o.TakerAddress.Hex(),
		"senderAddress":         o.SenderAddress.Hex(),
		"feeRecipientAddress":   o.FeeRecipientAddress.Hex(),
		"makerAssetData":        o.MakerAssetData,
		"makerFeeAssetData":     o.MakerFeeAssetData,
		"takerAssetData":        o.TakerAssetData,
		"takerFeeAssetData":     o.TakerFeeAssetData,
		"salt":                  o.Salt.String(),
		"makerFee":              o.MakerFee.String(),
		"takerFee":              o.TakerFee.String(),
		"makerAssetAmount":      o.MakerAssetAmount.String(),
		"takerAssetAmount":      o.TakerAssetAmount.String(),
		"expirationTimeSeconds": o.ExpirationTimeSeconds.String(),
	}

	var typedData = gethsigner.TypedData{
		Types:       eip712OrderTypes,
		PrimaryType: "Order",
		Domain:      domain,
		Message:     message,
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return common.Hash{}, err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return common.Hash{}, err
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hashBytes := keccak256(rawData)
	hash := common.BytesToHash(hashBytes)
	o.hash = &hash
	return hash, nil
}

// SignOrder signs the 0x order with the supplied Signer
func SignOrder(signer Signer, order *Order) (*SignedOrder, error) {
	if order == nil {
		return nil, errors.New("cannot sign nil order")
	}
	orderHash, err := order.ComputeOrderHash()
	if err != nil {
		return nil, err
	}

	ecSignature, err := signer.EthSign(orderHash.Bytes(), order.MakerAddress)
	if err != nil {
		return nil, err
	}

	// Generate 0x EthSign Signature (append the signature type byte)
	signature := make([]byte, 66)
	signature[0] = ecSignature.V
	copy(signature[1:33], ecSignature.R[:])
	copy(signature[33:65], ecSignature.S[:])
	signature[65] = byte(EthSignSignature)
	signedOrder := &SignedOrder{
		Order:     *order,
		Signature: signature,
	}
	return signedOrder, nil
}

// OrderJSON is an unmodified JSON representation of an Order
type OrderJSON struct {
	ChainID               int64  `json:"chainId"`
	ExchangeAddress       string `json:"exchangeAddress"`
	MakerAddress          string `json:"makerAddress"`
	MakerAssetData        string `json:"makerAssetData"`
	MakerFeeAssetData     string `json:"makerFeeAssetData"`
	MakerAssetAmount      string `json:"makerAssetAmount"`
	MakerFee              string `json:"makerFee"`
	TakerAddress          string `json:"takerAddress"`
	TakerAssetData        string `json:"takerAssetData"`
	TakerFeeAssetData     string `json:"takerFeeAssetData"`
	TakerAssetAmount      string `json:"takerAssetAmount"`
	TakerFee              string `json:"takerFee"`
	SenderAddress         string `json:"senderAddress"`
	FeeRecipientAddress   string `json:"feeRecipientAddress"`
	ExpirationTimeSeconds string `json:"expirationTimeSeconds"`
	Salt                  string `json:"salt"`
}

// SignedOrderJSON is an unmodified JSON representation of a SignedOrder
type SignedOrderJSON struct {
	ChainID               int64  `json:"chainId"`
	ExchangeAddress       string `json:"exchangeAddress"`
	MakerAddress          string `json:"makerAddress"`
	MakerAssetData        string `json:"makerAssetData"`
	MakerFeeAssetData     string `json:"makerFeeAssetData"`
	MakerAssetAmount      string `json:"makerAssetAmount"`
	MakerFee              string `json:"makerFee"`
	TakerAddress          string `json:"takerAddress"`
	TakerAssetData        string `json:"takerAssetData"`
	TakerFeeAssetData     string `json:"takerFeeAssetData"`
	TakerAssetAmount      string `json:"takerAssetAmount"`
	TakerFee              string `json:"takerFee"`
	SenderAddress         string `json:"senderAddress"`
	FeeRecipientAddress   string `json:"feeRecipientAddress"`
	ExpirationTimeSeconds string `json:"expirationTimeSeconds"`
	Salt                  string `json:"salt"`
	Signature             string `json:"signature"`
}

// UnmarshalJSON implements a custom JSON unmarshaller for the SignedOrder type
func (s *SignedOrder) UnmarshalJSON(data []byte) error {
	var signedOrderJSON SignedOrderJSON
	err := json.Unmarshal(data, &signedOrderJSON)
	if err != nil {
		return err
	}
	var ok bool
	s.ChainID = big.NewInt(signedOrderJSON.ChainID)
	s.ExchangeAddress = common.HexToAddress(signedOrderJSON.ExchangeAddress)
	s.MakerAddress = common.HexToAddress(signedOrderJSON.MakerAddress)
	s.MakerAssetData = common.FromHex(signedOrderJSON.MakerAssetData)
	s.MakerFeeAssetData = common.FromHex(signedOrderJSON.MakerFeeAssetData)
	if signedOrderJSON.MakerAssetAmount != "" {
		s.MakerAssetAmount, ok = math.ParseBig256(signedOrderJSON.MakerAssetAmount)
		if !ok {
			s.MakerAssetAmount = nil
		}
	}
	if signedOrderJSON.MakerFee != "" {
		s.MakerFee, ok = math.ParseBig256(signedOrderJSON.MakerFee)
		if !ok {
			s.MakerFee = nil
		}
	}
	s.TakerAddress = common.HexToAddress(signedOrderJSON.TakerAddress)
	s.TakerAssetData = common.FromHex(signedOrderJSON.TakerAssetData)
	s.TakerFeeAssetData = common.FromHex(signedOrderJSON.TakerFeeAssetData)
	if signedOrderJSON.TakerAssetAmount != "" {
		s.TakerAssetAmount, ok = math.ParseBig256(signedOrderJSON.TakerAssetAmount)
		if !ok {
			s.TakerAssetAmount = nil
		}
	}
	if signedOrderJSON.TakerFee != "" {
		s.TakerFee, ok = math.ParseBig256(signedOrderJSON.TakerFee)
		if !ok {
			s.TakerFee = nil
		}
	}
	s.SenderAddress = common.HexToAddress(signedOrderJSON.SenderAddress)
	s.FeeRecipientAddress = common.HexToAddress(signedOrderJSON.FeeRecipientAddress)
	if signedOrderJSON.ExpirationTimeSeconds != "" {
		s.ExpirationTimeSeconds, ok = math.ParseBig256(signedOrderJSON.ExpirationTimeSeconds)
		if !ok {
			s.ExpirationTimeSeconds = nil
		}
	}
	if signedOrderJSON.Salt != "" {
		s.Salt, ok = math.ParseBig256(signedOrderJSON.Salt)
		if !ok {
			s.Salt = nil
		}
	}
	s.Signature = common.FromHex(signedOrderJSON.Signature)
	return nil
}

// MarshalJSON implements a custom JSON marshaller for the SignedOrder type
func (s *SignedOrder) MarshalJSON() ([]byte, error) {
	makerAssetData := ""
	if len(s.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerAssetData))
	}
	makerFeeAssetData := "0x"
	if len(s.MakerFeeAssetData) != 0 {
		makerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerFeeAssetData))
	}
	takerAssetData := ""
	if len(s.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerAssetData))
	}
	takerFeeAssetData := "0x"
	if len(s.TakerFeeAssetData) != 0 {
		takerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerFeeAssetData))
	}
	signature := ""
	if len(s.Signature) != 0 {
		signature = fmt.Sprintf("0x%s", common.Bytes2Hex(s.Signature))
	}

	signedOrderBytes, err := json.Marshal(SignedOrderJSON{
		ChainID:               s.ChainID.Int64(),
		ExchangeAddress:       strings.ToLower(s.ExchangeAddress.Hex()),
		MakerAddress:          strings.ToLower(s.MakerAddress.Hex()),
		MakerAssetData:        makerAssetData,
		MakerFeeAssetData:     makerFeeAssetData,
		MakerAssetAmount:      s.MakerAssetAmount.String(),
		MakerFee:              s.MakerFee.String(),
		TakerAddress:          strings.ToLower(s.TakerAddress.Hex()),
		TakerAssetData:        takerAssetData,
		TakerFeeAssetData:     takerFeeAssetData,
		TakerAssetAmount:      s.TakerAssetAmount.String(),
		TakerFee:              s.TakerFee.String(),
		SenderAddress:         strings.ToLower(s.SenderAddress.Hex()),
		FeeRecipientAddress:   strings.ToLower(s.FeeRecipientAddress.Hex()),
		ExpirationTimeSeconds: s.ExpirationTimeSeconds.String(),
		Salt:                  s.Salt.String(),
		Signature:             signature,
	})
	return signedOrderBytes, err
}

// UnmarshalJSON implements a custom JSON unmarshaller for the SignedOrder type
func (s *Order) UnmarshalJSON(data []byte) error {
	var orderJson OrderJSON
	err := json.Unmarshal(data, &orderJson)
	if err != nil {
		return err
	}
	var ok bool
	s.ChainID = big.NewInt(orderJson.ChainID)
	s.ExchangeAddress = common.HexToAddress(orderJson.ExchangeAddress)
	s.MakerAddress = common.HexToAddress(orderJson.MakerAddress)
	s.MakerAssetData = common.FromHex(orderJson.MakerAssetData)
	s.MakerFeeAssetData = common.FromHex(orderJson.MakerFeeAssetData)
	if orderJson.MakerAssetAmount != "" {
		s.MakerAssetAmount, ok = math.ParseBig256(orderJson.MakerAssetAmount)
		if !ok {
			s.MakerAssetAmount = nil
		}
	}
	if orderJson.MakerFee != "" {
		s.MakerFee, ok = math.ParseBig256(orderJson.MakerFee)
		if !ok {
			s.MakerFee = nil
		}
	}
	s.TakerAddress = common.HexToAddress(orderJson.TakerAddress)
	s.TakerAssetData = common.FromHex(orderJson.TakerAssetData)
	s.TakerFeeAssetData = common.FromHex(orderJson.TakerFeeAssetData)
	if orderJson.TakerAssetAmount != "" {
		s.TakerAssetAmount, ok = math.ParseBig256(orderJson.TakerAssetAmount)
		if !ok {
			s.TakerAssetAmount = nil
		}
	}
	if orderJson.TakerFee != "" {
		s.TakerFee, ok = math.ParseBig256(orderJson.TakerFee)
		if !ok {
			s.TakerFee = nil
		}
	}
	s.SenderAddress = common.HexToAddress(orderJson.SenderAddress)
	s.FeeRecipientAddress = common.HexToAddress(orderJson.FeeRecipientAddress)
	if orderJson.ExpirationTimeSeconds != "" {
		s.ExpirationTimeSeconds, ok = math.ParseBig256(orderJson.ExpirationTimeSeconds)
		if !ok {
			s.ExpirationTimeSeconds = nil
		}
	}
	if orderJson.Salt != "" {
		s.Salt, ok = math.ParseBig256(orderJson.Salt)
		if !ok {
			s.Salt = nil
		}
	}
	return nil
}

// MarshalJSON implements a custom JSON marshaller for the Order type
func (s *Order) MarshalJSON() ([]byte, error) {
	makerAssetData := ""
	if len(s.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerAssetData))
	}
	makerFeeAssetData := "0x"
	if len(s.MakerFeeAssetData) != 0 {
		makerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerFeeAssetData))
	}
	takerAssetData := ""
	if len(s.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerAssetData))
	}
	takerFeeAssetData := "0x"
	if len(s.TakerFeeAssetData) != 0 {
		takerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerFeeAssetData))
	}

	orderBytes, err := json.Marshal(OrderJSON{
		ChainID:               s.ChainID.Int64(),
		ExchangeAddress:       strings.ToLower(s.ExchangeAddress.Hex()),
		MakerAddress:          strings.ToLower(s.MakerAddress.Hex()),
		MakerAssetData:        makerAssetData,
		MakerFeeAssetData:     makerFeeAssetData,
		MakerAssetAmount:      s.MakerAssetAmount.String(),
		MakerFee:              s.MakerFee.String(),
		TakerAddress:          strings.ToLower(s.TakerAddress.Hex()),
		TakerAssetData:        takerAssetData,
		TakerFeeAssetData:     takerFeeAssetData,
		TakerAssetAmount:      s.TakerAssetAmount.String(),
		TakerFee:              s.TakerFee.String(),
		SenderAddress:         strings.ToLower(s.SenderAddress.Hex()),
		FeeRecipientAddress:   strings.ToLower(s.FeeRecipientAddress.Hex()),
		ExpirationTimeSeconds: s.ExpirationTimeSeconds.String(),
		Salt:                  s.Salt.String(),
	})
	return orderBytes, err
}
