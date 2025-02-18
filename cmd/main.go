package main

import (
	"wallet/api/router"
	"wallet/config"
)

func main() {
	config.InitConfig()
	config.InitChainConfig()
	r := router.Router()
	err := r.Run(":8081")
	if err != nil {
		panic(err)
	}
}
