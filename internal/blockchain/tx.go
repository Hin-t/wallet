package blockchain

import (
	"context"
	"demo1/config"
	"demo1/internal/grpc"
	"demo1/internal/repository"
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
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
	"log"
	"strconv"
	"time"
)

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

// Send 发送交易
func Send(name, password, receiverStr, amountStr string) *tx.BroadcastTxResponse {
	privateKey, _ := wallet.LoadKey(name, password)

	// 构造交易
	// 初始化交易配置
	txConfig := initTxConfig()
	txBuilder := txConfig.NewTxBuilder()
	// 构造消息
	senderStr := wallet.GetAccountAddress(name)
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
	resp, err := Broadcast(txBytes)
	if err != nil {
		return nil
	}
	fmt.Println("响应消息：", resp.TxResponse)

	// 存储链下交易数据
	start := time.Now()
	repository.CreateOffChainTransaction(resp.TxResponse.TxHash, senderStr, receiverStr, config.WasmdChainID, amount)
	elapsed := time.Since(start)
	fmt.Println("time spend", elapsed)

	return resp
}

// MultiSend 多笔转账
func MultiSend(name, password string, receiversStr, amountsStr []string) *tx.BroadcastTxResponse {
	privateKey, _ := wallet.LoadKey(name, password)

	// 构造交易
	// 初始化交易配置
	txConfig := initTxConfig()
	txBuilder := txConfig.NewTxBuilder()
	// 构造消息
	senderStr := wallet.GetAccountAddress(name)

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
	resp, err := Broadcast(txBytes)
	if err != nil {
		log.Fatal("广播失败：", err)
	}
	fmt.Println("响应消息：", resp.TxResponse)

	// 存储链下交易数据
	start := time.Now()
	repository.CreateOffChainTransaction(resp.TxResponse.TxHash, senderStr, "multi-send", config.WasmdChainID, utils.FromString2Int64(amountsStr[0]))
	elapsed := time.Since(start)
	fmt.Println("time spend", elapsed)

	return resp
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
	sender := wallet.GetAccountAddress(name)

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
	repository.CreateOffChainTransaction(resp.TxResponse.TxHash, sender, recipient, config.WasmdChainID, parseInt)
	elapsed := time.Since(start)
	fmt.Println("time spend", elapsed)

	return resp
}

// DistributeCurrencyInRotation  发币,使用send逐个对用户账户发送金币
func DistributeCurrencyInRotation() {
	accounts := repository.GetAccountList()
	//ch := make(chan string, 100)
	amount := 1000/len(accounts) + 1
	for _, account := range accounts {
		Send("alice-test", "", account.Address, strconv.Itoa(amount))
		//WriteToRedis(utils.Rdb, ch, "alice-test", "", account.Address, "3")
	}
	//BroadcastToBlockchain(utils.Rdb, ch)
}

// DistributeCurrencyByMultiSend 使用MultiSend进行发币
func DistributeCurrencyByMultiSend() {
	accounts := repository.GetAccountList()
	addresses := repository.GetReceiversList()
	amounts := make([]string, len(accounts))
	for i := 0; i < len(accounts); i++ {
		amounts[i] = strconv.Itoa(1000/len(accounts) + 1)
	}
	MultiSend("alice-test", "", addresses, amounts)
}
