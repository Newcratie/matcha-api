package api

import (
	"fmt"
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

	var err error
	app.Neo, err = driver.OpenNeo("bolt://neo4j:secret@" + host + ":7687")
	fmt.Println(err)

	go app.routerAPI()
	app.Db = dbConnect()

	app.R.Run(":81")
}
