package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

func adminError(err error, c *gin.Context) {
	log.Error("login Error: ", err.Error())
	c.JSON(401, gin.H{
		"username": "",
		"token":    false,
	})
}

func checkAdmin(username string, c *gin.Context, Users []User) (User, error) {
	for _, user := range Users {
		if user.Username == username {
			return user, nil
		}
	}
	loginError(errors.New("username doesn't exist"), c)
	return User{}, errors.New("bad")
}
