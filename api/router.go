package api

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func (app *App) routerAPI() {
	m := melody.New()
	auth := app.R.Group("/auth")
	{
		auth.POST("/login", Login)
		auth.POST("/register", Register)
		auth.GET("/valid_email/:token", Token)
	}
	api := app.R.Group("/api")
	{
		api.POST("/add_like", CreateLike)
		api.POST("/get_people", GetPeople)
		api.POST("/get_matchs", GetMatchs)
		api.POST("/get_messages", GetMessages)
		api.POST("/user", UserHandler)
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
