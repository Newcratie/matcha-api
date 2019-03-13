package api

import (
	"fmt"
	"github.com/Newcratie/matcha-api/api/logprint"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

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

func tableOf(values string) string {
	return strings.Replace(values, ":", "", 99999)
}

func (app *App) insertUser(u User) {
	var query = `INSERT INTO public.users ` + tableOf(vUsers) + ` VALUES ` + vUsers
	_, err := app.Db.NamedExec(query, u)
	if err != nil {
		log.Fatalln(err)
	}
}
