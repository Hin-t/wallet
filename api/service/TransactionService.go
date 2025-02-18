package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wallet/internal/blockchain"
)

// SendToReceiver
// @Summary SendToReceiver
// @Description send sendAddress to receiverAddress amount
// @Tags 交易模块
// @param name formData string false "发送账户名称"
// @param password formData string false "密码"
// @param receiverAddress formData string false "接收地址"
// @param amount formData string false "转账金额"
// @Success 200 {string} json{"code","message"}
// @Router /transaction/send [post]
func SendToReceiver(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	receiverAddress := c.PostForm("receiverAddress")
	amount := c.PostForm("amount")
	resp := blockchain.Send(name, password, receiverAddress, amount)
	c.JSON(http.StatusOK, resp.TxResponse)

}

// QueryBalance
// @Summary QueryBalance
// @Description send sendAddress to receiverAddress amount
// @Tags 交易模块
// @param address query string false "查询地址"
// @Success 200 {string} json{"code","message"}
// @Router /transaction/queryBalance [get]
func QueryBalance(c *gin.Context) {
	address := c.Query("address")
	resp := blockchain.QueryBalance(address)
	c.JSON(http.StatusOK, gin.H{"balance": resp.Balance})
}

// QueryTransaction
// @Summary QueryTransaction
// @Description query transaction by txhash
// @Tags 交易模块
// @param txhash formData string false "交易哈希"
// @Success 200 {string} json{"code","message"}
// @Router /transaction/queryTransaction [post]
func QueryTransaction(c *gin.Context) {
	txhash := c.PostForm("txhash")
	resp := blockchain.QueryOnChainTransaction(txhash)
	c.JSON(http.StatusOK, gin.H{
		"query_tx_hash": txhash,
		"tx_hash":       resp.TxHash,
		"gas_used":      resp.GasUsed,
		"gas_wanted":    resp.GasWanted,
		"height":        resp.Height,
		"timestamp":     resp.Timestamp,
		"events":        resp.Events,
	})
}

// ContractQueryBalance
// @Summary ContractQueryBalance
// @Description query contract balance by address
// @Tags 交易模块
// @param query_addr formData string false "查询地址"
// @Success 200 {string} json{"code","message"}
// @Router /transaction/queryContractBalance [post]
func ContractQueryBalance(c *gin.Context) {
	QueryAddr := c.PostForm("query_addr")
	resp := blockchain.QueryBalanceByContract(QueryAddr)
	c.JSON(http.StatusOK, resp)
}

// TransferCW20TokenService
// @Summary TransferCW20TokenService
// @Description transfer by contract
// @Tags 交易模块
// @param name formData string false "发送账户名称"
// @param password formData string false "密码"
// @param recipient formData string false "接收地址"
// @param amount formData string false "转账金额"
// @Success 200 {string} json{"code","message"}
// @Router /transaction/transferCW20TokenService [post]
func TransferCW20TokenService(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	recipient := c.PostForm("recipient")
	amount := c.PostForm("amount")
	resp := blockchain.TransferCW20Token(name, password, recipient, amount)
	c.JSON(http.StatusOK, resp)
}

// QueryTokenInfo
// @Summary QueryTokenInfo
// @Description query token info
// @Tags 交易模块
// @Success 200 {string} json{"code","message"}
// @Router /transaction/queryTokenInfo [get]
func QueryTokenInfo(c *gin.Context) {
	resp := blockchain.QueryTokenInfo()
	c.JSON(http.StatusOK, resp)
}
