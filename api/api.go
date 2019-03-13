package api

import (
	"github.com/gin-gonic/gin"
)

const hashKey = "5c894d411c2f7445dbadb9b6"

func (app *App) newApp() {
	app.R = gin.Default()
}

func Run() {
	app.newApp()
	go app.routerAPI()
	app.Db = dbConnect()

	app.R.Run(":81")
}
