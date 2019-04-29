package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"time"
)

func (app *App) routerAPI() {
	m := melody.New()
	app.R.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
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
		api.POST("/add_like", CreateLike)
		api.GET("/people", GetPeople)
		api.GET("/matchs", GetMatchs)
		api.GET("/messages", GetMessages)
		api.POST("/user", UserHandler)
		api.PUT("/user/:name", UserModifyHandler)
		api.POST("/img/:n", UserImageHandler)
		api.GET("/ws/:user/:suitor", func(c *gin.Context) {
			_ = m.HandleRequest(c.Writer, c.Request)
		})
	}
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		app.insertMessage(msg)
		_ = m.BroadcastFilter(msg, func(session *melody.Session) bool {
			//AUth: verify if token is valid here.
			return session.Request.URL.Path == s.Request.URL.Path
		})
	})
}
