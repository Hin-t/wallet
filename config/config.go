package config

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
)

// Config 结构体存储所有配置信息
type Config struct {
	RPC   RPCConfig   `mapstructure:"rpc"`
	Redis RedisConfig `mapstructure:"redis"`
	Log   LogConfig   `mapstructure:"log"`
}

// RPCConfig 存储区块链节点配置信息
type RPCConfig struct {
	CosmosEndpoint string `mapstructure:"cosmos_endpoint"`
	GRPCEndpoint   string `mapstructure:"grpc_endpoint"`
}

// RedisConfig 存储 Redis 相关配置
type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LogConfig 存储日志级别
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// GlobalConfig 变量存储加载后的配置
var GlobalConfig Config

// LoadConfig 读取 YAML 配置文件
func LoadConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析到结构体
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	fmt.Println("✅ 配置加载成功")
	return nil
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	fmt.Println("config: ", viper.Get("config"))
}

func InitChainConfig() {
	// 获取 SDK 全局配置
	config := sdk.GetConfig()

	// 设置新的地址前缀（例如 "mychain"）
	config.SetBech32PrefixForAccount("wasm", "wasm")

	// 也可以设置验证人和共识节点的前缀
	//config.SetBech32PrefixForValidator("mychainvaloper", "mychainvaloperpub")
	//config.SetBech32PrefixForConsensusNode("mychainvalcons", "mychainvalconspub")

	// 确保配置生效
	config.Seal()
}
