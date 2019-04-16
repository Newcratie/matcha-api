package api

import (
	"fmt"
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
		api.POST("/get_matchs", GetMatchs)
		api.GET("/ws/:token/:chan", func(c *gin.Context) {
			m.HandleRequest(c.Writer, c.Request)
		})
	}

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Println("MSG ===> ", string(msg))
		m.BroadcastFilter(msg, func(session *melody.Session) bool {
			//AUth: verify if token is valid here.
			fmt.Println("s: ", s.Request.URL.Path)
			fmt.Println("session: ", session.Request.URL.Path)
			return session.Request.URL.Path == s.Request.URL.Path
		})
	})
}
