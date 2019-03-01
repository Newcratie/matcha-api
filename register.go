package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
	"strings"
	"time"
)

const vUsers = `(:username, :email, :lastname, :firstname, :password, :random_token, :img1, :img2, :img3, :img4, :img5, :biography, :birthday, :genre, :interest, :city, :zip, :country, :latitude, :longitude, :geo_allowed, :online, :rating, :admin)`

func Register(c *gin.Context) {
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
	user, err := validateUser(rf)
	if err != nil {
		c.JSON(401, gin.H{"err": err.Error()})
	} else {
		fmt.Println("register success", user)
		c.JSON(401, gin.H{"err": "good"})
	}
	//app.insertUser(user)
	//c.JSON(200, user)
}

func (app *App) checkRegister(rf registerForm) (error) {
	// here rf should not exist on DB, password must match confirm, validity of all datas.
	u := User{}
	err := app.Db.Get(&u, `SELECT * FROM users WHERE username=$1 OR email=$2`, rf.Username, rf.Email)
	fmt.Println("checkRegister: >>", u.Id, "<<", err)
	if u.Id != 0 {
		return errors.New("Username or Email already exist")
	}
	return nil
}

func parseTime(str string) (time.Time, error) {
	var re = regexp.MustCompile(`\s\((.*)\)`)
	s := re.ReplaceAllString(str, ``)
	t, err := time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", s)
	return t, err
}

func NewUser(rf registerForm) User {
	birthday, _ := parseTime(rf.Birthday)
	u := User{0,
		rf.Username,
		rf.Email,
		rf.Lastname,
		rf.Firstname,
		rf.Password,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		birthday,
		"",
		"",
		"",
		"",
		"",
		0,
		0,
		false,
		false,
		0,
		false,
		"",
	}
	return u
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
