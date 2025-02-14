package repository

import (
	"demo1/internal/db"
	"demo1/internal/wallet"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

type Account struct {
	Name      string
	Publickey string
	Address   string
	Mnemonic  string //助记词存储方式，字符串？字符串列表
}

func (table *Account) TableName() string {
	return "Account"
}

// CreateAccount 使用keyring为指定账户名称生成公私钥和地址
func CreateAccount(name, password string) *Account {
	record, mnemonic := wallet.CreateNewAccount(name, password)
	address, _ := record.GetAddress()
	account := &Account{
		Name:      record.Name,
		Address:   address.String(),
		Publickey: record.PubKey.String(),
		Mnemonic:  mnemonic,
	}
	db.InitMySQL().Create(account)
	return account
}

// GetAccountList 获取数据库中全部账户信息
func GetAccountList() []*Account {
	data := make([]*Account, 5)
	db.InitMySQL().Find(&data)
	return data
}

// GetAccountAddressByName 根据给定账户名称查找用户地址
func GetAccountAddressByName(name string) string {
	account := &Account{}
	db.InitMySQL().Where("name = ?", name).First(account)
	return account.Address
}

// GetAccountInfoByName 获取本地指定名称的密钥环记录
func GetAccountInfoByName(name string) *keyring.Record {
	kr := wallet.InitKeyring()
	record, err := kr.Key(name)
	if err != nil {
		panic(err)
	}
	return record
}

// GetReceiversList 获取接收者列表
func GetReceiversList() []string {
	accounts := GetAccountList()
	addresses := make([]string, 0)
	for _, account := range accounts {
		addresses = append(addresses, account.Address)
	}
	return addresses
}
