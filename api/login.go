package api

import (
	"github.com/Newcratie/matcha-api/api/hash"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

func loginError(err error, c *gin.Context) {
	log.Error("login Error: ", err.Error())
	c.JSON(441, gin.H{
		"username": "",
	})
}
func (app *App) getUser(Username string) (u User, err error) {
	err = app.Db.Get(&u, `SELECT * FROM users WHERE username=$1`, Username)
	return
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	u, err := app.getUser(username)
	if err != nil {
		c.JSON(401, gin.H{"err": "User doesn't exist"})
	} else if password != hash.Decrypt(hashKey, u.Password) {
		c.JSON(401, gin.H{"err": "Wrong password"})
	} else {
		c.JSON(200, u)
	}
}
