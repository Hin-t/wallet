package blockchain

import (
	"context"
	"demo1/config"
	"demo1/internal/grpc"
	"demo1/internal/repository"
	"encoding/json"
	"fmt"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"log"
	"sync"
)

var Wg sync.WaitGroup

type TokenInfoResponse struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    uint8  `json:"decimals"`
	TotalSupply string `json:"total_supply"`
}

// QueryBalanceByContract 调用智能合约查询账户余额
func QueryBalanceByContract(QueryAddr string) *wasmtypes.QuerySmartContractStateResponse {
	//创建Wasm查询客户端
	wasmClient := wasmtypes.NewQueryClient(grpc.NewGrpcConn())

	//生成查询JSON
	queryMsg := map[string]interface{}{
		"balance": map[string]string{"address": QueryAddr},
	}
	queryBytes, _ := json.Marshal(queryMsg)

	// 发送查询请求
	res, err := wasmClient.SmartContractState(
		context.Background(),
		&wasmtypes.QuerySmartContractStateRequest{
			Address:   config.CONTRACT,
			QueryData: queryBytes,
		})
	if err != nil {
		log.Fatalf("查询 CW20 余额失败: %v", err)
	}
	// 解析返回的 JSON 结果
	var balanceRes map[string]string
	if err := json.Unmarshal(res.Data, &balanceRes); err != nil {
		log.Fatalf("解析 JSON 失败: %v", err)
	}

	// 输出余额
	fmt.Printf("账户 %s 的 CW20 余额: %s\n", QueryAddr, balanceRes["balance"])
	return res
}

// QueryTokenInfo 查询 CW20 代币信息
func QueryTokenInfo() *wasmtypes.QuerySmartContractStateResponse {
	wasmClient := wasmtypes.NewQueryClient(grpc.NewGrpcConn())

	tokenInfoMsg := map[string]interface{}{
		"token_info": struct{}{},
	}
	tokenInfoBytes, _ := json.Marshal(tokenInfoMsg)
	// 发送查询请求
	res, err := wasmClient.SmartContractState(
		context.Background(),
		&wasmtypes.QuerySmartContractStateRequest{
			Address:   config.CONTRACT,
			QueryData: tokenInfoBytes,
		})
	if err != nil {
		log.Fatalf("查询 token_info 失败: %v", err)
	}
	// 解析返回的 JSON 结果
	var tokenInfo TokenInfoResponse
	if err := json.Unmarshal(res.Data, &tokenInfo); err != nil {
		log.Fatalf("解析 JSON 失败: %v", err)
	}
	fmt.Printf("Token Name: %s, Symbol: %s, Decimals: %d, Total Supply: %s\n",
		tokenInfo.Name, tokenInfo.Symbol, tokenInfo.Decimals, tokenInfo.TotalSupply)
	return res
}

// DistributeCW0Currency  合约发币
func DistributeCW0Currency() {
	defer Wg.Done()
	accounts := repository.GetAccountList()
	for _, account := range accounts {
		TransferCW20Token("alice-test", "", account.Address, "1")
	}
}
