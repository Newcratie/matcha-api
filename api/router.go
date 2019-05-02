package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"time"
)

func (app *App) routerAPI() {
	app.M = melody.New()
	app.R.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := app.R.Group("/auth")
	{
		auth.POST("/login", Login)
		auth.POST("/register", Register)
		auth.GET("/valid_email/:token", Token)
	}
	api := app.R.Group("/api")
	{
		api.GET("/people", GetPeople)
		api.PUT("/people/:id/:action", CreateLike)
		//Here handle action: like, dislike or block
		//then return the same thing than GetPeople please
		api.GET("/matchs", GetMatchs)
		api.GET("/messages", GetMessages)
		api.GET("/user", UserHandler)
		api.PUT("/user/:name", UserModify)
		api.POST("/img/:n", UserImageHandler)
		api.GET("/notifications/history/:user", notificationsHistoryHandler)
		api.GET("/notifications/websocket/:user", func(c *gin.Context) {
			_ = app.M.HandleRequest(c.Writer, c.Request)
		})
		api.GET("/ws/:user/:suitor", func(c *gin.Context) {
			_ = app.M.HandleRequest(c.Writer, c.Request)
		})
	}
	app.M.HandleMessage(func(s *melody.Session, msg []byte) {
		app.insertMessage(msg)
		_ = app.M.BroadcastFilter(msg, func(session *melody.Session) bool {
			//AUth: verify if token is valid here.
			return session.Request.URL.Path == s.Request.URL.Path
		})
	})

	for i := 0; i < 3000; i++ {
		time.Sleep(time.Second * 2)
		fmt.Println("i++")
		n := Notification{
			"Ceci est une notifications test",
			int64(i),
			100,
			23,
			45,
		}
		msg, _ := json.Marshal(n)
		app.postNotification(n, msg)
	}
}
