package api

import (
	"encoding/json"
	"fmt"
	"github.com/Newcratie/matcha-api/api/hash"
	"github.com/Newcratie/matcha-api/api/logprint"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	InfoC    = "\033[1;34m%s\033[0m"
	NoticeC  = "\033[1;36m%s\033[0m"
	WarningC = "\033[1;33m%s\033[0m"
	ErrorC   = "\033[1;31m%s\033[0m"
	DebugC   = "\033[0;36m%s\033[0m"
)

func Token(c *gin.Context) {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User{random_token : "`+c.Param("token")+`"}) SET n.access_lvl = 1 RETURN n`, nil)
	if len(data) == 0 {
		c.JSON(201, gin.H{"err": "Wrong token"})
	} else if data[0][0].(graph.Node).Properties["access_lvl"] == int64(1) {
		c.JSON(201, gin.H{"err": "Email already validated"})
	} else {
		c.JSON(200, gin.H{"status": "Email validated"})
	}
}

func PrintHandlerLog(Err string, Color string) {
	Err = Err + "\n"
	fmt.Printf(Color, Err)
}

func CreateLike(c *gin.Context) {

	claims := jwt.MapClaims{}
	valid, err := ValidateToken(c, &claims)

	id := int(claims["id"].(float64))
	UpdateLastConn(id)

	if valid == true {
		var m Match
		m.idTo, _ = strconv.Atoi(c.Param("id"))
		m.action = strings.ToUpper(c.Param("action"))
		m.idFrom = id
		m.c = c
		if _, err = app.dbMatchs(m); err != nil {
			c.JSON(201, gin.H{"err": err.Error()})
		} else {
			c.JSON(200, nil)
		}
		app.onlineRefresh(strconv.Itoa(m.idFrom))

	} else {
		PrintHandlerLog("Token Not Valid", ErrorC)
		fmt.Println("jwt error: ", err)
		c.JSON(201, gin.H{"err": err.Error()})
	}
}

func ValidateToken(c *gin.Context, claims jwt.Claims) (valid bool, err error) {
	tokenString := c.Request.Header["Authorization"][0]

	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(hashKey), nil
	})
	if err != nil {
		fmt.Println("jwt error: ", err)
		c.JSON(201, gin.H{"err": err.Error()})
		return false, err
	} else if checkJwt(tokenString) {
		return true, err
	}
	return false, err
}

func Next(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(200, gin.H{
		"next": "next",
	})
}

func GetMatchs(c *gin.Context) {
	tokenString := c.Request.Header["Authorization"][0]

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(hashKey), nil
	})

	id := int(math.Round(claims["id"].(float64)))
	UpdateLastConn(id)

	if err != nil {
		c.JSON(202, gin.H{"err": err.Error()})
	} else if checkJwt(tokenString) {
		app.onlineRefresh(strconv.Itoa(id))
		g, err := app.dbGetMatchs(id)
		if err != nil {
			c.JSON(201, gin.H{"err": err.Error()})
		} else {
			c.JSON(200, g)
		}
	}
}
func GetMessages(c *gin.Context) {
	tokenString := c.Request.Header["Authorization"][0]
	suitorId := c.Request.Header["Suitor-Id"][0]

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(hashKey), nil
	})

	id := int(math.Round(claims["id"].(float64)))
	UpdateLastConn(id)

	if err != nil {
		c.JSON(202, gin.H{"err": err.Error()})
	} else if checkJwt(tokenString) {
		sui, _ := strconv.Atoi(suitorId)
		app.onlineRefresh(strconv.Itoa(id))
		msgs, err := app.dbGetMessages(id, sui)
		if err != nil {
			c.JSON(201, gin.H{"err": err.Error()})
		} else {
			c.JSON(200, msgs)
		}
	}
}

func GetPeople(c *gin.Context) {
	filtersJson := c.Request.Header["Filters"][0]
	var err error

	filters := Filters{}
	claims := jwt.MapClaims{}

	valid, err := ValidateToken(c, &claims)
	json.Unmarshal([]byte(filtersJson), &filters)

	id := int(claims["id"].(float64))

	if err != nil {
		c.JSON(202, gin.H{"err": err.Error()})
	} else if valid == true {
		str := strconv.Itoa(id)
		app.onlineRefresh(str)
		app.alertOnline(true, str)
		g, err := app.dbGetPeople(id, &filters)
		if err != nil {
			c.JSON(201, gin.H{"err": err.Error()})
		} else {
			UpdateLastConn(id)
			c.JSON(200, g)
		}
	} else {
		fmt.Println("jwt error: ", err)
		c.JSON(201, gin.H{"err": err.Error()})
	}
}

func newVisit(c *gin.Context) {
	newEvent(c, func(name string) string {
		return name + " has visited your profil page"
	})
	c.JSON(200, gin.H{})
}

func UpdateLastConn(Id int) {
	u, _ := app.getUser(Id, "")
	u.Id = int64(Id)
	u.LastConn, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339Nano))
	u.Online = true
	app.updateLastConn(u)
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	u, err := app.getUser(-1, username)
	if err != nil || password != hash.Decrypt(hashKey, u.Password) {
		c.JSON(201, gin.H{"err": "Wrong password or username"})
	} else if u.AccessLvl == 0 {
		c.JSON(201, gin.H{"err": "You must validate your Email first"})
	} else {
		jwt, err := u.GenerateJwt()
		if err != nil {
			c.JSON(201, gin.H{"err": "Internal server error: " + err.Error()})
		} else {
			c.JSON(200, jwt)
		}
	}
}

func Register(c *gin.Context) {
	logprint.Title("Register")
	fmt.Println("POST BIRTHDAY =========", c.PostForm("birthday"), "|")
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
		c.JSON(201, res)
	} else {
		app.insertUser(user)
		c.JSON(200, gin.H{})
	}
	logprint.End()
}
