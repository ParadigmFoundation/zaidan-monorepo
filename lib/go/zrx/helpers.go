package zrx

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"

	"github.com/0xProject/0x-mesh/ethereum/wrappers"

	"github.com/0xProject/0x-mesh/ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var NULL_ADDRESS = common.Address{}

type contracts struct {
	exchange    *wrappers.Exchange
	exchangeABI abi.ABI

	devUtils *wrappers.DevUtilsCaller
}

type ZeroExHelper struct {
	ChainID           *big.Int
	ContractAddresses ethereum.ContractAddresses

	client         *ethclient.Client
	contracts      *contracts
	orderValidator *ordervalidator.OrderValidator
}

func NewZeroExHelper(client *ethclient.Client, maxContentLength int) (*ZeroExHelper, error) {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	addresses, err := ethereum.GetContractAddressesForChainID(int(chainID.Int64()))
	if err != nil {
		return nil, err
	}

	exchangeABI, err := abi.JSON(strings.NewReader(wrappers.ExchangeABI))
	if err != nil {
		return nil, err
	}

	exchange, err := wrappers.NewExchange(addresses.Exchange, client)
	if err != nil {
		return nil, err
	}

	devUtils, err := wrappers.NewDevUtilsCaller(addresses.DevUtils, client)
	if err != nil {
		return nil, err
	}

	orderValidator, err := ordervalidator.New(client, int(chainID.Int64()), maxContentLength)
	if err != nil {
		return nil, err
	}

	return &ZeroExHelper{
		ChainID:           chainID,
		ContractAddresses: addresses,
		client:            client,
		orderValidator:    orderValidator,
		contracts:         &contracts{exchange: exchange, exchangeABI: exchangeABI, devUtils: devUtils},
	}, nil
}

// DevUtils returns an initialized 0x DevUtils contract caller
func (zh *ZeroExHelper) DevUtils() *wrappers.DevUtilsCaller { return zh.contracts.devUtils }

func (zh *ZeroExHelper) OrderValidator() *ordervalidator.OrderValidator { return zh.orderValidator }

// GetFillOrderCallData generates the underlying 0x exchange call data for the fill (to be singed by the taker)
func (zh *ZeroExHelper) GetFillOrderCallData(order zeroex.Order, takerAssetAmount *big.Int, signature []byte) ([]byte, error) {
	return zh.contracts.exchangeABI.Pack("fillOrder", order, takerAssetAmount, signature)
}

// GetTransactionHash gets the 0x transaction hash for the current chain ID
func (zh *ZeroExHelper) GetTransactionHash(tx *Transaction) (common.Hash, error) {
	return tx.ComputeHashForChainID(int(zh.ChainID.Int64()))
}

// ValidateFill is a convenience wrapper for ordervalidator.BatchValidate with a single order
func (zh *ZeroExHelper) ValidateFill(ctx context.Context, order *zeroex.SignedOrder, takerAssetAmount *big.Int) error {
	orders := []*zeroex.SignedOrder{order}
	rawValidationResults := zh.orderValidator.BatchValidate(ctx, orders, true, rpc.LatestBlockNumber)

	if len(rawValidationResults.Rejected) == 1 {
		return fmt.Errorf("%s", rawValidationResults.Rejected[0].Status.Message)
	}

	if len(rawValidationResults.Accepted) != 1 {
		return fmt.Errorf("unable to validate order")
	}

	// if taker is null address, skip taker checks
	nullAddress := common.Address{}
	if order.TakerAddress == nullAddress {
		return nil
	}

	takerBalanceInfo, err := zh.contracts.devUtils.GetBalanceAndAssetProxyAllowance(nil, order.TakerAddress, order.TakerAssetData)
	if err != nil {
		return err
	}

	if order.TakerAssetAmount.Cmp(takerBalanceInfo.Allowance) > 0 {
		return fmt.Errorf("taker has insufficient allowance for trade: (has: %s), (want: %s)", takerBalanceInfo.Allowance, order.TakerAssetAmount)
	}
	if order.TakerAssetAmount.Cmp(takerBalanceInfo.Balance) > 0 {
		return fmt.Errorf("taker has insufficient balance for trade: (has: %s), (want: %s)", takerBalanceInfo.Balance, order.TakerAssetAmount)
	}

	return nil
}

func (zh *ZeroExHelper) CreateOrder(
	maker common.Address,
	taker common.Address,
	sender common.Address,
	feeRecipient common.Address,
	makerAsset common.Address,
	takerAsset common.Address,
	makerAmount *big.Int,
	takerAmount *big.Int,
	makerFee *big.Int,
	takerFee *big.Int,
	makerFeeAsset common.Address,
	takerFeeAsset common.Address,
	expirationTimeSeconds *big.Int,
) (*zeroex.Order, error) {
	salt, err := GeneratePseudoRandomSalt()
	if err != nil {
		return nil, err
	}

	return &zeroex.Order{
		ChainID:               zh.ChainID,
		ExchangeAddress:       zh.ContractAddresses.Exchange,
		MakerAddress:          maker,
		MakerAssetData:        EncodeERC20AssetData(makerAsset),
		MakerFeeAssetData:     EncodeERC20AssetData(makerFeeAsset),
		MakerAssetAmount:      makerAmount,
		MakerFee:              makerFee,
		TakerAddress:          taker,
		TakerAssetData:        EncodeERC20AssetData(takerAsset),
		TakerFeeAssetData:     EncodeERC20AssetData(takerFeeAsset),
		TakerAssetAmount:      takerAmount,
		TakerFee:              takerFee,
		SenderAddress:         sender,
		FeeRecipientAddress:   feeRecipient,
		ExpirationTimeSeconds: expirationTimeSeconds,
		Salt:                  salt,
	}, nil
}
