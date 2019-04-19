package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func (app *App) insertMessage(byt []byte) {
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	dat["author"] = int(dat["author"].(float64))
	dat["to"] = int(dat["to"].(float64))
	q := `MATCH (a:User),(b:User)
WHERE ID(a)={author} AND ID(b)={to}
CREATE (a)-[s:SAYS]->(message:Message {msg:{msg}, author: {author}, to:{id}, timestamp:{timestamp}})-[t:TO]->(b)`
	st := app.prepareStatement(q)
	executeStatement(st, dat)
}

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
		api.GET("/ws/:user/:suitor", func(c *gin.Context) {
			m.HandleRequest(c.Writer, c.Request)
		})
	}

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(session *melody.Session) bool {
			//AUth: verify if token is valid here.
			fmt.Printf("MSG ===> %T\n", msg)
			fmt.Println("s: ", s.Request.URL.Path)
			fmt.Println("session: ", session.Request.URL.Path)
			app.insertMessage(msg)
			return session.Request.URL.Path == s.Request.URL.Path
		})
	})
}
