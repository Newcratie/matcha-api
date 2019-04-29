package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Newcratie/matcha-api/api/hash"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
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

func getBodymap(c *gin.Context) (body map[string]interface{}) {
	r, _ := c.GetRawData()
	err := json.Unmarshal(r, &body)
	if err != nil {
		panic(err)
	}
	return
}

func UserImageHandler(c *gin.Context) {
	//
	file := c.PostForm("file")
	fmt.Println("file  ===>", file)
	//
	//claims := jwt.MapClaims{}
	//valid, err := ValidateToken(c, &claims)
	//if valid {
	//	Id := int(claims["id"].(float64))
	//	g, err := app.dbGetUserProfile(Id)
	//	tagList := app.dbGetTagList()
	//	if err != nil {
	//		c.JSON(201, gin.H{"err": err.Error()})
	//	} else {
	//		c.JSON(200, gin.H{"user": g, "tagList": tagList})
	//	}
	//} else {
	//	c.JSON(201, gin.H{"err": err.Error()})
	//}
}

func UserModifyHandler(c *gin.Context) {
	if strings.Contains(c.Param("name"), "img") {
		file, err := c.FormFile("file")
		fmt.Println("file  ===>", file, err)
	} else {
		m := getBodymap(c)
		fmt.Println("Map  ===>", m)
	}
	claims := jwt.MapClaims{}
	valid, err := ValidateToken(c, &claims)
	if valid {
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
	fmt.Println("ON Pass change")
	fmt.Println("Claims ==>", claims)

	username := claims["username"].(string)
	mail := claims["email"].(string)
	//oldPassword := c.PostForm("old_password")
	//newPassword := c.PostForm("new_password")
	//confirmPassword := c.PostForm("confirm_password")

	oldPassword := "123456789"
	newPassword := "Pouet1234/"
	confirmPassword := "Pouet1234/"

	u, err := app.getUser(username)
	if err != nil || oldPassword != hash.Decrypt(hashKey, u.Password) {
		fmt.Println("Wrong Pass")
		err = errors.New("error : wrong password")
	} else {
		err = verifyPassword(newPassword, confirmPassword)
		if err != nil {
			fmt.Println("ERROR :===> ", err)
			c.JSON(201, gin.H{"err": err.Error()})
		} else {
			fmt.Println("Password change::")
			u.Password = hash.Encrypt(hashKey, newPassword)
			app.updateUser(u)
			err = SendEmail("Matcha password change", username, mail, "./api/utils/pass_change.html")
		}
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

		case "tag":

		case "password":
			UserPassChange(c, claims)
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
