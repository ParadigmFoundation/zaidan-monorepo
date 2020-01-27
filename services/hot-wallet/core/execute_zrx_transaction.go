package core

import (
	"context"
	"math/big"
	"time"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/zrx"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ParadigmFoundation/zaidan-monorepo/lib/go/grpc"
)

// ExecuteZeroExTransaction implements grpc.HotWalletServer
func (hw *HotWallet) ExecuteZeroExTransaction(ctx context.Context, req *grpc.ExecuteZeroExTransactionRequest) (*grpc.ExecuteZeroExTransactionResponse, error) {
	ztx, err := req.Transaction.ToZeroExTransaction()
	if err != nil {
		return nil, err
	}

	nonce, err := hw.provider.Client().PendingNonceAt(ctx, hw.senderAddress)
	if err != nil {
		return nil, err
	}

	opts := &bind.TransactOpts{
		From:     hw.senderAddress,
		Signer:   hw.senderTransactor.Signer,
		Nonce:    new(big.Int).SetUint64(nonce),
		Value:    new(big.Int).Mul(zrx.PROTOCOL_FEE_MULTIPLIER, ztx.GasPrice),
		GasPrice: new(big.Int).Set(ztx.GasPrice),
		GasLimit: zrx.EXECUTE_FILL_TX_GAS_LIMIT,
		Context:  ctx,
	}

	tx, err := hw.zrxHelper.ExecuteTransaction(opts, ztx, req.Signature)
	if err != nil {
		return nil, err
	}

	return &grpc.ExecuteZeroExTransactionResponse{
		TransactionHash: tx.Hash().Hex(),
		SubmittedAt:     time.Now().UnixNano() / 1e6, // convert from NS to MS
	}, nil
}
