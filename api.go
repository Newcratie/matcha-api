package api

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

func (app *App) newApp() {
	app.R = gin.Default()
	app.Users = make([]User, 0)
}

func Run() {
	app.newApp()
	go app.routerAPI()
	app.Db = dbConnect()
	app.fetchUsers()
	app.validate = validator.New()
	app.R.Run(":81")
}
