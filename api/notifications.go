package api

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"strconv"
)

func (app *App) postNotification(n Notification, msg []byte) {
	id := strconv.FormatInt(n.UserId, 10)
	url := "/api/notifications/websocket/" + id
	app.M.BroadcastFilter(msg, func(session *melody.Session) bool {
		return session.Request.URL.Path == url
	})
}
func notificationsHistoryHandler(c *gin.Context) {
	n := []Notification{}
	c.JSON(200, n)
}
