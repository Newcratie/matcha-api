package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"os"
	"time"
)

const hashKey = "5c894d411c2f7445dbadb9b6"

var prin = fmt.Println

func (app *App) newApp() {
	app.R = gin.Default()
}

func NewConn(host string) (bolt.Conn, error) {
	var err error
	app.Db, err = bolt.NewDriverPool("bolt://neo4j:secret@"+host+":7687", 1000)
	if err != nil {
		return nil, err
	}
	retries := 0
	for retries < 300 {
		conn, _ := app.Db.OpenPool()
		if conn != nil {
			return conn, nil
		}
		retries = retries + 1
		fmt.Println("neo4j not ready, waiting 15s and trying again ", retries)
		time.Sleep(15 * time.Second)
	}
	return nil, err
}

func Run() {
	app.newApp()
	host := os.Getenv("NEO_HOST")
	app.Neo, _ = NewConn(host)
	defer app.Neo.Close()

	go app.routerAPI()
	app.R.Run(":81")
}
