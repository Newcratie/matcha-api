package api

import (
	"github.com/Newcratie/matcha-api/api/kwal"
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
		api.PUT("/visit/:user_id", newVisit)
		api.PUT("/people/:id/:action", CreateLike)
		api.GET("/matchs", GetMatchs)
		api.GET("/kwal", func(c *gin.Context) {
			k := kwal.GetKeys()
			c.JSON(200, k)
		})
		api.GET("/messages", GetMessages)
		api.GET("/user", UserHandler)
		api.PUT("/user/:name", UserModify)
		api.POST("/img/:n", userImageHandler)
		api.GET("/notifications/history/:user", notificationsHistoryHandler)
		api.DELETE("/notifications/:id", notificationsDeleteHandler)
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
}
