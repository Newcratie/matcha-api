package api

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

func loginError(err error, c *gin.Context) {
	log.Error("login Error: ", err.Error())
	c.JSON(441, gin.H{
		"username": "",
	})
}
