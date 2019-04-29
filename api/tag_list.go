package api

import (
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"strconv"
)

func (app *App) dbGetTagList() (ret []Tag) {
	q := `
MATCH (a:TAG)
RETURN a
ORDER BY a.value
`
	data, _, _, _ := app.Neo.QueryNeoAll(q, map[string]interface{}{})
	for _, tab := range data {
		ret = append(ret, Tag{
			tab[0].(graph.Node).Properties["key"].(string),
			tab[0].(graph.Node).Properties["text"].(string),
			tab[0].(graph.Node).Properties["value"].(string),
		})
	}
	return
}

func (app *App) insertTag(t Tag, Id int64) {
	fmt.Println("========", MapOf(t))
	q := `MATCH (u:User) WHERE ID(u) = ` + strconv.Itoa(Id) + ` CREATE (t:TAG{key: {key}, text:{text}, value:{value}})<-[:TAGGED]-(u)`
	st := app.prepareStatement(q)
	executeStatement(st, MapOf(t))
}
