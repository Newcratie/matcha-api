package api

import (
	"errors"
	"fmt"
	"github.com/Newcratie/matcha-api/api/hash"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UserHandler(c *gin.Context) {
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

func UserPassChange(c *gin.Context, claims jwt.MapClaims) {
	username := claims["username"].(string)
	mail := claims["mail"].(string)
	password := c.PostForm("password")

	u, err := app.getUser(username)
	if err != nil || password != hash.Decrypt(hashKey, u.Password) {
		err = errors.New("Err: Wrong Password")
	} else {
		// check password
		//change password
		err = SendEmail(username, mail, "./api/utils/pass_change.html")
	}
}

func UserModify(c *gin.Context) {
	claims := jwt.MapClaims{}
	valid, err := ValidateToken(c, &claims)
	if valid == false || err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
	} else {
		mod := c.Param("name")
		Id := int(claims["id"].(float64))
		switch mod {
		case "biography":
			updateBio(c, Id)
		case "username":

		case "password":

		case "firstname":

		case "lastname":

		case "location":
		}
	}
}

func updateBio(c *gin.Context, Id int) {
	fmt.Println("TATATATA")
}

func updatePassword(c *gin.Context, Id int) {

}

func updateUsername(c *gin.Context, Id int) {

}

func updateFirstname(c *gin.Context, Id int) {

}

func updateLastname(c *gin.Context, Id int) {

}

func updateLocation(c *gin.Context, Id int) {

}
