package wallet

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
)

// CreateNewAccount 生成新的账户
func CreateNewAccount(name, password string) (*keyring.Record, string) {
	kr := InitKeyring()
	// 创建账户
	accountName := name
	/*
		此处需要判断新建账户名称是否存在
	*/
	record, mnemonic, err := kr.NewMnemonic(accountName, keyring.English, sdk.FullFundraiserPath, password, hd.Secp256k1)
	if err != nil {
		log.Fatalf("Failed to create account: %v\n", err)
	}
	return record, mnemonic
}

// ImportAccount 导入账户（通过助记词）
func ImportAccount(name, mnemonic, password string) *keyring.Record {
	kr := InitKeyring()

	record, err := kr.NewAccount(name, mnemonic, password, sdk.FullFundraiserPath, hd.Secp256k1)
	if err != nil {
		log.Fatalf("Failed to create account: %v\n", err)
	}
	return record
}

// GetAccount 获取账户信息
func GetAccount(name string) *keyring.Record {
	kr := InitKeyring()
	record, err := kr.Key(name)
	if err != nil {
		log.Fatalf("Failed to get account: %v\n", err)
	}
	return record
}

// GetAccountAddress 获取账户地址
func GetAccountAddress(name string) string {
	record := GetAccount(name)
	addr, _ := record.GetAddress()
	return addr.String()
}
