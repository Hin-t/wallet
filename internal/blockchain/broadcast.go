// Package blockchain
/*
	交易广播
*/
package blockchain

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"wallet/internal/grpc"
)

// Broadcast 广播交易
func Broadcast(txBytes []byte) (*tx.BroadcastTxResponse, error) {
	txClient := tx.NewServiceClient(grpc.NewGrpcConn())
	resp, err := txClient.BroadcastTx(context.Background(), &tx.BroadcastTxRequest{
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
		TxBytes: txBytes,
	})
	return resp, err
}
