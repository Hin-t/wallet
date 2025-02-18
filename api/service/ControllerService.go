package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wallet/internal/repository"
)

func LoginService(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	c.JSON(http.StatusOK, gin.H{"code": repository.Login(name, password)})
}

func RegisterService(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	c.JSON(http.StatusOK, gin.H{"success": repository.Register(name, password)})
}
