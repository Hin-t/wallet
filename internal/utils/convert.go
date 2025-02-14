package utils

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"strconv"
)

// InitCodec 初始化编码器
func InitCodec() *codec.ProtoCodec {
	interfaceRegistry := types.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry) // 注册加密相关类型
	cdc := codec.NewProtoCodec(interfaceRegistry)
	return cdc
}

// FromString2Int64 字符串转int64
func FromString2Int64(str string) int64 {
	parseInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return parseInt
}
