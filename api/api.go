package api

import (
	"github.com/gin-gonic/gin"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"os"
)

const hashKey = "5c894d411c2f7445dbadb9b6"

func (app *App) newApp() {
	app.R = gin.Default()
}

func Run() {
	app.newApp()
	driver := bolt.NewDriver()
	host := os.Getenv("NEO_HOST")
	//host := "localhost"
	app.Neo, _ = driver.OpenNeo("bolt://neo4j:secret@" + host + ":7687")

	go app.routerAPI()
	app.Db = dbConnect()

	app.R.Run(":81")
}
