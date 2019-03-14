package api

import (
	"fmt"
	"github.com/Newcratie/matcha-api/api/hash"
	"github.com/Newcratie/matcha-api/api/logprint"
	"github.com/gin-gonic/gin"
	"time"
)

func Token(c *gin.Context) {
	logprint.Title("Token Handler")
	err := app.validToken(c.Param("token"))
	if err != nil {
		c.JSON(401, gin.H{"err": err.Error()})
	} else {
		c.JSON(200, gin.H{"status": "Email validated"})
	}
	logprint.End()
}

func Next(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(200, gin.H{
		"next": "next",
	})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	u, err := app.getUser(username)
	if err != nil {
		c.JSON(401, gin.H{"err": "User doesn't exist"})
	} else if password != hash.Decrypt(hashKey, u.Password) {
		c.JSON(401, gin.H{"err": "Wrong password"})
	} else if u.AccessLvl == 0 {
		c.JSON(401, gin.H{"err": "You must validate your Email first"})
	} else {
		c.JSON(200, u)
	}
}

func Register(c *gin.Context) {
	logprint.Title("Register")
	bd, _ := time.Parse(time.RFC3339, c.PostForm("birthday"))
	rf := registerForm{
		c.PostForm("username"),
		c.PostForm("email"),
		c.PostForm("password"),
		c.PostForm("confirm"),
		c.PostForm("lastname"),
		c.PostForm("firstname"),
		bd,
	}
	user, res := validateUser(rf)
	if !res.Valid {
		c.JSON(401, res)
	} else {
		fmt.Println("register success", user)
		app.insertUser(user)
		c.JSON(200, gin.H{})
	}
	logprint.End()
}
