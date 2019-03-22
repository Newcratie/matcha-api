package api

import (
	"fmt"
	"github.com/Newcratie/matcha-api/api/hash"
	"github.com/Newcratie/matcha-api/api/logprint"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"math"
	"time"
)

func Token(c *gin.Context) {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User{random_token : "`+c.Param("token")+`"}) SET n.access_lvl = 1 RETURN n`, nil)
	if len(data) == 0 {
		c.JSON(401, gin.H{"err": "Wrong token"})
	} else if data[0][0].(graph.Node).Properties["access_lvl"] == int64(1) {
		c.JSON(401, gin.H{"err": "Email already validated"})
	} else {
		c.JSON(200, gin.H{"status": "Email validated"})
	}
}

func Next(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(200, gin.H{
		"next": "next",
	})
}

func Start(c *gin.Context) {
	tokenString := c.Request.Header["Authorization"][0]

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(hashKey), nil
	})
	fmt.Println(claims)
	if err != nil {
		fmt.Println(err)
		c.JSON(401, gin.H{"err": err.Error()})
	} else if checkJwt(tokenString) {
		id := int(math.Round(claims["id"].(float64)))
		u, err := app.getBasicUser(id)
		fmt.Println(u)
		if err != nil {
			fmt.Println(err)
			c.JSON(401, gin.H{"err": err.Error()})
		} else {
			d, err := app.getBasicDates(id)
			if err != nil {
				fmt.Println(err)
				c.JSON(401, gin.H{"err": err.Error()})
			} else {
				c.JSON(200, ResStart{u, d})
			}
		}
	}
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	u, err := app.getUser(username)
	fmt.Println(u, err)
	if err != nil || password != hash.Decrypt(hashKey, u.Password) {
		c.JSON(401, gin.H{"err": "Wrong password or username"})
	} else if u.AccessLvl == 0 {
		c.JSON(401, gin.H{"err": "You must validate your Email first"})
	} else {
		jwt, err := u.GenerateJwt()
		if err != nil {
			c.JSON(401, gin.H{"err": "Internal server error: " + err.Error()})
		} else {
			c.JSON(200, jwt)
		}
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
