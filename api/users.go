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
	file := c.PostForm("file")
	fmt.Printf("file  %s\n", file)

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

func UserModify(c *gin.Context) {
	claims := jwt.MapClaims{}
	valid, err := ValidateToken(c, &claims)
	if valid == false || err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
	} else {
		mod := c.Param("name")
		switch mod {
		case "biography":
			updateBio(c, claims)
			break
		case "username":
			updateUsername(c, claims)
			break
		case "tag":
			addTag(c, claims)
			break
		case "password":
			updatePassword(c, claims)
			break
		case "firstname":
			updateFirstname(c, claims)
			break
		case "lastname":
			updateLastname(c, claims)
			break
			//case "location":
		}
	}
}

func updateBio(c *gin.Context, claims jwt.MapClaims) {
	username := claims["username"].(string)
	u, err := app.getUser(username)
	if err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
		return
	}

	bio := c.PostForm("bio")
	if len(bio) > 100 {
		err = errors.New("error : your biography can't exceed 100 characters")
		c.JSON(201, gin.H{"err": err.Error()})
	} else {
		u.Biography = bio
		app.updateUser(u)
	}
}

func updateUsername(c *gin.Context, claims jwt.MapClaims) {
	username := claims["username"].(string)
	u, err := app.getUser(username)
	if err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
		return
	}

	newUsername := c.PostForm("username")
	if len(newUsername) < 6 || len(newUsername) > 20 {
		err = errors.New("error : your username must be between 6 to 20 characters")
		c.JSON(201, gin.H{"err": err.Error()})
	} else {
		u.Username = newUsername
		app.updateUser(u)
	}
}

func addTag(c *gin.Context, claims jwt.MapClaims) {
	var Tags Tag
	username := claims["username"].(string)
	u, err := app.getUser(username)
	if err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
		return
	}

	Tags.Value = c.PostForm("tag")
	if len(Tags.Value) < 1 || len(Tags.Value) > 20 {
		err = errors.New("error : your Tag must be between 1 to 20 characters")
		c.JSON(201, gin.H{"err": err.Error()})
	} else {
		Tags.Value = strings.ToLower(Tags.Value)
		Tags.Key = Tags.Value
		Tags.Text = "#" + strings.Title(Tags.Value)
		app.insertTag(Tags, u.Id)
	}
}

func updatePassword(c *gin.Context, claims jwt.MapClaims) {
	fmt.Println("ON Pass change")
	fmt.Println("Claims ==>", claims)

	username := claims["username"].(string)
	mail := claims["email"].(string)
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	confirmPassword := c.PostForm("confirm_password")

	//oldPassword := "123456789"
	//newPassword := "Pouet1234/"
	//confirmPassword := "Pouet1234/"

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

func updateFirstname(c *gin.Context, claims jwt.MapClaims) {
	username := claims["username"].(string)
	u, err := app.getUser(username)
	if err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
		return
	}

	firstname := c.PostForm("firstname")
	if len(firstname) < 2 || len(firstname) > 20 {
		err = errors.New("error : your firstname must be between 2 to 20 characters")
		c.JSON(201, gin.H{"err": err.Error()})
	} else {
		u.FirstName = firstname
		app.updateUser(u)
	}
}

func updateLastname(c *gin.Context, claims jwt.MapClaims) {
	username := claims["username"].(string)
	u, err := app.getUser(username)
	if err != nil {
		c.JSON(201, gin.H{"err": err.Error()})
		return
	}

	lastname := c.PostForm("lastname")
	if len(lastname) < 2 || len(lastname) > 20 {
		err = errors.New("error : your lastname must be between 2 to 20 characters")
		c.JSON(201, gin.H{"err": err.Error()})
	} else {
		u.LastName = lastname
		app.updateUser(u)
	}
}

//check kat long validity
//func updateLocation(c *gin.Context, claims jwt.MapClaims) {
//	username := claims["username"].(string)
//	u, err := app.getUser(username)
//	if err != nil {
//		c.JSON(201, gin.H{"err": err.Error()})
//		return
//	}
//
//	lat := c.PostForm("latitude")
//	lon := c.PostForm("longitude")
//	if len(firstname) < 2 {
//		err = errors.New("error : your firstname must be at least 2 characters")
//		c.JSON(201, gin.H{"err": err.Error()})
//	} else {
//		u.FirstName = firstname
//		app.updateUser(u)
//	}
//}
