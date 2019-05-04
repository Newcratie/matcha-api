package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"strconv"
)

type Notification struct {
	Message   string `json:"message"`
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	AuthorId  int64  `json:"author_id"`
	SubjectId int64  `json:"subject_id"`
}

func (app *App) postNotification(message string, userId, authorId, subjectId int64) {
	n := Notification{
		message,
		0,
		userId,
		authorId,
		subjectId,
	}
	msg, _ := json.Marshal(n)
	id := strconv.FormatInt(34, 10)
	url := "/api/notifications/websocket/" + id

	_ = app.M.BroadcastFilter(msg, func(session *melody.Session) bool {
		return session.Request.URL.Path == url
	})
}

func (app *App) dbInsertNotification(byt []byte) {
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
CREATE (a)<-[s:TO]-(n:Notif {message:{message}, author_id: {author_id}, subject_id: {subject_id}})`
	st := app.prepareStatement(q)
	executeStatement(st, dat)
}

func notificationsHistoryHandler(c *gin.Context) {
	n := []Notification{}
	c.JSON(200, n)
}
