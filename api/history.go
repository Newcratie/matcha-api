package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"strconv"
)

type History struct {
	Message   string `json:"message"`
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	AuthorId  int64  `json:"author_id"`
	SubjectId int64  `json:"subject_id"`
}

func getHistoryHandler(c *gin.Context) {
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

