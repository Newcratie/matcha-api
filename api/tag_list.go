package api

import (
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
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

func (app *App) insertTag(t Tag) {
	fmt.Println("========", MapOf(t))
	q := `CREATE (t:Tag{key: {key}, text:{text}, value:{value}})`
	st := app.prepareStatement(q)
	executeStatement(st, MapOf(t))
}
