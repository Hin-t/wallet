package router

import (
	"demo1/api/service"
	"demo1/docs"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	router := gin.Default()
	// swagger
	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 创建账户
	router.POST("/user/createAccount", service.CreateAccount)
	// 查找所有账户信息
	router.GET("/user/findAllAccount", service.FindAllAccount)
	// 查找指定账户地址
	router.POST("/user/findAccountAddressByName", service.FindAccountAddressByName)
	// 发起单笔转账交易 A -> B
	router.POST("/transaction/send", service.SendToReceiver)
	// 查询address的账户余额
	router.POST("/transaction/queryBalance", service.QueryBalance)
	// 查询txhash的的交易详情
	router.POST("/transaction/queryTransaction", service.QueryTransaction)
	// 查询智能合约地址余额
	router.POST("/transaction/queryContractBalance", service.ContractQueryBalance)
	// 调用智能合约进行转账
	router.POST("/transaction/transferCW20TokenService", service.TransferCW20TokenService)
	// 查询token_info
	router.GET("/transaction/queryTokenInfo", service.QueryTokenInfo)

	return router
}
