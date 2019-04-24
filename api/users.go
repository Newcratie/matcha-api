package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UserHandler(c *gin.Context) {
	//i := "Changed"

	claims := jwt.MapClaims{}
	valid, err := ValidateToken(c, &claims)

	if valid == true {
		Id := int(claims["id"].(float64))
		g, err := app.dbGetUserProfile(Id)
		tagList := app.dbGetTagList()
		if err != nil {
			c.JSON(201, gin.H{"err": err.Error()})
		} else {
			c.JSON(200, gin.H{"user": g, "tagList": tagList})
		}
	} else {
		c.JSON(201, gin.H{"err": err.Error()})
	}
}

func UserPassMailChange(Id int, claims jwt.MapClaims) {

	if claims["password"] != "" {

	}
}
