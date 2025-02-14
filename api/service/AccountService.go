package service

import (
	"demo1/internal/repository"
	"demo1/internal/wallet"
	"github.com/gin-gonic/gin"
)

// CreateAccount
// @Summary CreateAccount
// @Description 创建账户
// @Tags 账户模块
// @param name formData string false "账户名称"
// @param password formData string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createAccount [Post]
func CreateAccount(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	user, _ := wallet.CreateNewAccount(name, password)
	c.JSON(200, gin.H{
		"code":    200,
		"message": user,
	})
}

// FindAllAccount
// @Summary FindAllAccount
// @Description 查找全部账户信息
// @Tags 账户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/findAllAccount [get]
func FindAllAccount(c *gin.Context) {
	data := repository.GetAccountList()
	c.JSON(200, data)
}

// FindAccountAddressByName
// @Summary FindAccountAddressByName
// @Description 根据给定账户名称查找账户地址
// @Tags 账户模块
// @param name formData string false "账户名称"
// @Success 200 {string} json{"code","message"}
// @Router /user/findAccountAddressByName [post]
func FindAccountAddressByName(c *gin.Context) {
	name := c.PostForm("name")
	address := repository.GetAccountAddressByName(name)
	c.JSON(200, gin.H{"address": address})
}
