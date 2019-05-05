package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"gopkg.in/olahol/melody.v1"
	"strconv"
)

func (app *App) postNotification(message string, userId, authorId, subjectId int64) {
	n := Notification{
		message,
		0,
		userId,
		authorId,
		subjectId,
	}
	id := strconv.FormatInt(userId, 10)
	url := "/api/notifications/websocket/" + id
	msg, _ := json.Marshal(n)
	n.Id = app.dbInsertNotification(msg)
	msg, _ = json.Marshal(n)

	_ = app.M.BroadcastFilter(msg, func(session *melody.Session) bool {
		return session.Request.URL.Path == url
	})
}

func (app *App) dbInsertNotification(byt []byte) int64 {
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	dat["message"] = dat["message"].(string)
	dat["author_id"] = int64(dat["author_id"].(float64))
	dat["user_id"] = int64(dat["user_id"].(float64))
	dat["subject_id"] = int64(dat["subject_id"].(float64))
	q := `
MATCH (a:User)
WHERE ID(a)={user_id}
CREATE (a)<-[s:TO]-(n:Notif {message:{message}, author_id: {author_id}, subject_id: {subject_id}})
RETURN ID(n)
`
	data, _, _, _ := app.Neo.QueryNeoAll(q, dat)
	return data[0][0].(int64)
}

func notificationsHistoryHandler(c *gin.Context) {
	n, _ := strconv.Atoi(c.Param("user"))
	userId := int64(n)
	fmt.Println(userId)
	q := `
MATCH (n:Notif)-[:TO]-(u:User) WHERE ID(u) = {user_id} RETURN n ORDER by ID(n)
`
	ntfs := make([]Notification, 0)
	data, _, _, _ := app.Neo.QueryNeoAll(q, map[string]interface{}{"user_id": userId})
	for _, tab := range data {
		ntfs = append(ntfs, Notification{
			tab[0].(graph.Node).Properties["message"].(string),
			int64(tab[0].(graph.Node).NodeIdentity),
			0,
			tab[0].(graph.Node).Properties["author_id"].(int64),
			tab[0].(graph.Node).Properties["subject_id"].(int64),
		})
	}
	c.JSON(200, ntfs)
}
func notificationsDeleteHandler(c *gin.Context) {
	q := `MATCH (n:Notif)-[r]-(u) WHERE ID(n) = ` + c.Param("id") + ` DELETE r, n`
	st := app.prepareStatement(q)
	executeStatement(st, map[string]interface{}{})
	c.JSON(200, nil)
}
