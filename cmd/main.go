package main

import (
	"demo1/api/router"
	"demo1/config"
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
