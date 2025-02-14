// Package blockchain
/*
	查询账户、交易状态*
*/
package blockchain

import (
	"context"
	"demo1/internal/grpc"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
	"log"
)

// QueryBalance 查询账户余额
func QueryBalance(address string) *sdk.Coin {
	// 新建查询客户端
	queryClient := types2.NewQueryClient(grpc.NewGrpcConn())
	req := &types2.QueryBalanceRequest{
		Address: address,
		Denom:   "stake",
	}
	resp, err := queryClient.Balance(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	return resp.Balance
}

// QueryAccountNumberAndSequence 查询账户的AccountNumber和Sequence
func QueryAccountNumberAndSequence(address string) (uint64, uint64) {
	queryClient := types.NewQueryClient(grpc.NewGrpcConn())
	req := &types.QueryAccountInfoRequest{
		Address: address,
	}
	info, err := queryClient.AccountInfo(context.Background(), req)
	if err != nil {
		panic(err)
	}
	return info.Info.AccountNumber, info.Info.Sequence
}

// QueryOnChainTransaction 查询链上指定交易
func QueryOnChainTransaction(txhash string) *sdk.TxResponse {
	serviceClient := tx.NewServiceClient(grpc.NewGrpcConn())
	req := &tx.GetTxRequest{
		Hash: txhash,
	}
	resp, err := serviceClient.GetTx(context.Background(), req)
	if err != nil {
		panic(err)
	}
	return resp.GetTxResponse()
}
