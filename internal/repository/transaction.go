package repository

import (
	"context"
	"demo1/config"
	"demo1/internal/db"
	"demo1/internal/grpc"
	"demo1/internal/utils"
	"demo1/internal/wallet"
	"encoding/json"
	"fmt"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	tx2 "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
	"log"
	"strconv"
	"sync"
	"time"
)

var Wg sync.WaitGroup

type Transaction struct {
	TxHash         string `gorm:"column:tx_hash"`                  // 链上交易哈希
	Sender         string `gorm:"column:sender"`                   // 发送方地址
	Receiver       string `gorm:"column:receiver"`                 // 接收方地址
	Amount         int64  `gorm:"column:amount"`                   // 交易金额
	Fee            int64  `gorm:"column:fee"`                      // 手续费
	Status         string `gorm:"column:status;default:'pending'"` // 状态（pending，confirmed，failed，cancelled）
	TxTimestamp    string `gorm:"column:timestamp"`                // 链上交易创建时间
	StoreTimestamp string `gorm:"column:store_timestamp"`          // 存储时间戳
	//signed_data     json.RawMessage gorm:"column:signed_data"     //签名的原始交易数据
	SignedData     string `gorm:"column:signed_data"`     // 签名的原始交易数据
	ChainReference string `gorm:"column:chain_reference"` // 对应链上的表示（链名或ID）
}

type TokenInfoResponse struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    uint8  `json:"decimals"`
	TotalSupply string `json:"total_supply"`
}

func (table *Transaction) TableName() string {
	return "transactions"
}

// CreateOffChainTransaction 创建链下交易记录
func CreateOffChainTransaction(txhash, sender, receiver, chainReference string, amount int64) {
	resp := pollQueryOnChainTransaction(txhash)
	/*
		判断 gas_limited 和 gas_wanted 大小关系 判断交易是否成功
	*/
	tx := &Transaction{
		TxHash:         txhash,
		Sender:         sender,
		Receiver:       receiver,
		Amount:         amount,
		Fee:            resp.GasUsed,
		Status:         "pending",
		TxTimestamp:    resp.Timestamp,
		StoreTimestamp: time.Now().String(),
		SignedData:     resp.Data,
		ChainReference: chainReference,
	}
	mysqldb := db.InitMySQL()
	mysqldb.Create(tx)
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

// PollQueryOnChainTransaction 查询链上指定交易
func pollQueryOnChainTransaction(txhash string) *sdk.TxResponse {
	serviceClient := tx.NewServiceClient(grpc.NewGrpcConn())
	req := &tx.GetTxRequest{
		Hash: txhash,
	}
	resp, err := serviceClient.GetTx(context.Background(), req)
	// 10秒内，查到交易后返回相应，否则返回nil
	for _ = range 200 {
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			resp, err = serviceClient.GetTx(context.Background(), req)
		} else {
			return resp.TxResponse
		}
	}
	return nil
}

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

// 初始化交易配置
func initTxConfig() client.TxConfig {
	// 启用的签名模式
	signModes := []signingtypes.SignMode{
		signingtypes.SignMode_SIGN_MODE_DIRECT,
	}
	return authtx.NewTxConfig(
		utils.InitCodec(), // Protobuf 编解码器
		signModes,         // 启用签名模式
	)
}

// constructSignerData 构造签名者信息
func constructSignerData(accountNumber, sequence uint64) xauthsigning.SignerData {
	// 构造签名者信息
	signerData := xauthsigning.SignerData{
		ChainID:       config.WasmdChainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
	return signerData
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

// GetAccountAddress 获取账户地址
func GetAccountAddress(name string) string {
	record := wallet.GetAccount(name)
	addr, _ := record.GetAddress()
	return addr.String()
}

// Send 发送交易
func Send(name, password, receiverStr, amountStr string) *tx.BroadcastTxResponse {
	privateKey, _ := wallet.LoadKey(name, password)

	// 构造交易
	// 初始化交易配置
	txConfig := initTxConfig()
	txBuilder := txConfig.NewTxBuilder()
	// 构造消息
	senderStr := GetAccountAddress(name)
	sender, err := sdk.AccAddressFromBech32(senderStr)
	if err != nil {
		log.Fatal("sender", err)
	}
	receiver, err := sdk.AccAddressFromBech32(receiverStr)
	if err != nil {
		log.Fatal("receiver", err)
	}
	amount := utils.FromString2Int64(amountStr)
	msg := types2.NewMsgSend(sender, receiver, sdk.NewCoins(sdk.NewInt64Coin("stake", amount)))
	err = txBuilder.SetMsgs(msg)
	if err != nil {
		panic(err)
	}
	// 设置Gas
	txBuilder.SetGasLimit(500000)
	// 构造签名
	accountNumber, sequence := QueryAccountNumberAndSequence(senderStr)

	var sigV2 signingtypes.SignatureV2
	sigV2 = signingtypes.SignatureV2{
		PubKey: privateKey.PubKey(),
		Data: &signingtypes.SingleSignatureData{
			SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		panic(err)
	}
	signerData := constructSignerData(accountNumber, sequence)
	sigV2, err = tx2.SignWithPrivKey(context.Background(),
		signingtypes.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder, privateKey, txConfig, sequence)
	if err != nil {
		panic(err)
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		panic(err)
	}
	// 产生Protobuf-encoded bytes.
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		panic(err)
	}

	/*
		将交易信息写进redis
		之后慢慢从redis中读取数据并广播
	*/

	// 广播交易
	txClient := tx.NewServiceClient(grpc.NewGrpcConn())
	resp, err := txClient.BroadcastTx(context.Background(), &tx.BroadcastTxRequest{
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
		TxBytes: txBytes,
	})
	if err != nil {
		log.Fatal("广播失败：", err)
	}
	fmt.Println("响应消息：", resp.TxResponse)

	// 存储链下交易数据
	start := time.Now()
	CreateOffChainTransaction(resp.TxResponse.TxHash, senderStr, receiverStr, config.WasmdChainID, amount)
	elapsed := time.Since(start)
	fmt.Println("time spend", elapsed)

	return resp
}

// getReceiversList获取接收者列表
func getReceiversList() []string {
	accounts := GetAccountList()
	addresses := make([]string, 0)
	for _, account := range accounts {
		addresses = append(addresses, account.Address)
	}
	return addresses
}

// MultiSend 多笔转账
func MultiSend(name, password string, receiversStr, amountsStr []string) *tx.BroadcastTxResponse {
	privateKey, _ := wallet.LoadKey(name, password)

	// 构造交易
	// 初始化交易配置
	txConfig := initTxConfig()
	txBuilder := txConfig.NewTxBuilder()
	// 构造消息
	senderStr := GetAccountAddress(name)

	if len(receiversStr) != len(amountsStr) {
		log.Fatal("接收者地址数量与金额数量不匹配")
	}

	// 构造 Inputs 和 Outputs
	var totalAmount sdk.Coins
	var inputs []types2.Input
	var outputs []types2.Output

	for i, receiver := range receiversStr {
		amount := utils.FromString2Int64(amountsStr[i])
		coin := sdk.NewInt64Coin("stake", amount)
		totalAmount = totalAmount.Add(coin)

		outputs = append(outputs, types2.Output{
			Address: receiver,
			Coins:   sdk.NewCoins(coin),
		})
	}

	inputs = append(inputs, types2.Input{
		Address: senderStr,
		Coins:   totalAmount,
	})

	msg := &types2.MsgMultiSend{
		Inputs:  inputs,
		Outputs: outputs,
	}

	err := txBuilder.SetMsgs(msg)
	if err != nil {
		panic(err)
	}
	// 设置Gas
	txBuilder.SetGasLimit(500000)
	// 构造签名
	accountNumber, sequence := QueryAccountNumberAndSequence(senderStr)

	var sigV2 signingtypes.SignatureV2
	sigV2 = signingtypes.SignatureV2{
		PubKey: privateKey.PubKey(),
		Data: &signingtypes.SingleSignatureData{
			SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		panic(err)
	}
	signerData := constructSignerData(accountNumber, sequence)
	sigV2, err = tx2.SignWithPrivKey(context.Background(),
		signingtypes.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder, privateKey, txConfig, sequence)
	if err != nil {
		panic(err)
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		panic(err)
	}
	// 产生Protobuf-encoded bytes.
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		panic(err)
	}

	/*
		将交易信息写进redis
		之后慢慢从redis中读取数据并广播
	*/

	// 广播交易
	txClient := tx.NewServiceClient(grpc.NewGrpcConn())
	resp, err := txClient.BroadcastTx(context.Background(), &tx.BroadcastTxRequest{
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
		TxBytes: txBytes,
	})
	if err != nil {
		log.Fatal("广播失败：", err)
	}
	fmt.Println("响应消息：", resp.TxResponse)

	// 存储链下交易数据
	//start := time.Now()
	//CreateOffChainTransaction(resp.TxResponse.TxHash, senderStr, receiverStr, utils.WasmdChainID, amount)
	//elapsed := time.Since(start)
	//fmt.Println("time spend", elapsed)

	return resp
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

// TransferCW20Token 使用智能合约进行转账
func TransferCW20Token(name, password, recipient, amount string) *tx.BroadcastTxResponse {
	// 创建 Wasm 客户端
	privateKey, _ := wallet.LoadKey(name, password)

	// 构造交易
	// 初始化交易配置
	txConfig := initTxConfig()
	txBuilder := txConfig.NewTxBuilder()
	// 构造消息
	sender := GetAccountAddress(name)

	transferMsg := map[string]interface{}{
		"transfer": map[string]string{
			"recipient": recipient,
			"amount":    amount,
		},
	}
	transferMsgBytes, err := json.Marshal(transferMsg)
	if err != nil {
		panic(err)
	}
	msg := &wasmtypes.MsgExecuteContract{
		Sender:   sender,
		Contract: config.CONTRACT,
		Msg:      transferMsgBytes,
	}
	err = txBuilder.SetMsgs(msg)
	if err != nil {
		panic(err)
	}
	// 设置Gas
	txBuilder.SetGasLimit(500000)
	// 构造签名
	accountNumber, sequence := QueryAccountNumberAndSequence(sender)

	var sigV2 signingtypes.SignatureV2
	sigV2 = signingtypes.SignatureV2{
		PubKey: privateKey.PubKey(),
		Data: &signingtypes.SingleSignatureData{
			SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		panic(err)
	}
	signerData := constructSignerData(accountNumber, sequence)
	sigV2, err = tx2.SignWithPrivKey(context.Background(),
		signingtypes.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder, privateKey, txConfig, sequence)
	if err != nil {
		panic(err)
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		panic(err)
	}
	// 产生Protobuf-encoded bytes.
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		panic(err)
	}
	// 广播交易
	txClient := tx.NewServiceClient(grpc.NewGrpcConn())
	resp, err := txClient.BroadcastTx(context.Background(), &tx.BroadcastTxRequest{
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
		TxBytes: txBytes,
	})
	if err != nil {
		log.Fatal("广播失败：", err)
	}
	fmt.Println("响应消息：", resp.TxResponse)

	parseInt, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return nil
	}
	// 存储链下交易数据
	start := time.Now()
	CreateOffChainTransaction(resp.TxResponse.TxHash, sender, recipient, config.WasmdChainID, parseInt)
	elapsed := time.Since(start)
	fmt.Println("time spend", elapsed)

	return resp
}
