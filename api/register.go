package api

import (
	"fmt"
	"github.com/Newcratie/matcha-api/api/logprint"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
	"strings"
	"time"
)

func Register(c *gin.Context) {
	logprint.Title("Register")
	admin := false
	if c.PostForm("admin") == "true" {
		admin = true
	}
	rf := registerForm{
		c.PostForm("username"),
		c.PostForm("email"),
		c.PostForm("password"),
		c.PostForm("confirm"),
		c.PostForm("lastname"),
		c.PostForm("firstname"),
		c.PostForm("birthday"),
		admin,
	}
	user, res := validateUser(rf)
	if res.Valid != true {
		logprint.PrettyPrint(res)
		c.JSON(401, res)
	} else {
		fmt.Println("register success", user)
		c.JSON(200, gin.H{})
	}
	app.insertUser(user)
	logprint.End()
}

func parseTime(str string) (time.Time, error) {
	var re = regexp.MustCompile(`\s\((.*)\)`)
	s := re.ReplaceAllString(str, ``)
	t, err := time.Parse(timeLayout, s)
	return t, err
}

func tableOf(values string) string {
	return strings.Replace(values, ":", "", 99999)
}

func (app *App) insertUser(u User) {
	var query = `INSERT INTO public.users ` + tableOf(vUsers) + ` VALUES ` + vUsers
	_, err := app.Db.NamedExec(query, u)
	if err != nil {
		log.Fatalln(err)
	} else {
		app.Users = append(app.Users, u)
	}
}
