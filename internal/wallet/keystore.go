package wallet

import (
	"demo1/internal/utils"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

// InitKeyring 初始化秘钥环
func InitKeyring() keyring.Keyring {
	// 配置keyring
	cdc := utils.InitCodec()

	keyringBackend := viper.GetString("keyring.backend")
	keyringDir := filepath.Join(".", viper.GetString("keyring.dir"))
	fmt.Println(keyringDir)
	kr, _ := keyring.New("cosmos", keyringBackend, keyringDir, os.Stdin, cdc)
	return kr
}

// LoadKey 加载私钥
func LoadKey(name string, password string) (types.PrivKey, string) {
	kr := InitKeyring()
	privateKeyArmor, err := kr.ExportPrivKeyArmor(name, password)
	if err != nil {
		log.Fatal("导出私钥失败：", err)
	}
	privateKey, algo, err := crypto.UnarmorDecryptPrivKey(privateKeyArmor, password)
	if err != nil {
		log.Fatal(err)
	}
	return privateKey, algo
}
