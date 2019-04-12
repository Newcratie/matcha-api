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
		api.POST("/get_people", GetPeople)
		api.GET("/ws", func(c *gin.Context) {
			m.HandleRequest(c.Writer, c.Request)
		})

	}

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})
}
