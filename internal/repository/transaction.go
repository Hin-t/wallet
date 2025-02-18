package repository

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"sync"
	"time"
	"wallet/internal/db"
	"wallet/internal/grpc"
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
