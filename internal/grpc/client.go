package grpc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
)

// NewGrpcConn 初始化Grpc连接
func NewGrpcConn() *grpc.ClientConn {
	// 建立连接
	conn, err := grpc.NewClient(viper.GetString("grpc.cosmwasm.endpoint"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	return conn
}
