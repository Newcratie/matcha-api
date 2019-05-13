package api

import (
	"encoding/json"
	"errors"
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"os"
	"strconv"
	"testing"
)

func getMatchingTest(t *testing.T) {
	host := os.Getenv("NEO_HOST")
	app.Db, _ = bolt.NewDriverPool("bolt://neo4j:secret@"+host+":7687", 1000)
	app.Neo, _ = app.Db.OpenPool()
	defer app.Neo.Close()
	page := 1
	skip := page * 50
	Id := 1315
	u, _ := getUser(Id, "")
	g := dbGetMatchingTest(u, skip)
	fmt.Println("G == ", g)
}

//MATCH (n:User), (u:User), (t:TAG) WHERE ID(n) = 1315 AND (t.value = n.tags[1]) AND NOT ID(u) = 1315 AND (u)-[:TAGGED]->(t) return u

func dbGetMatchingTest(u User, skip int) []graph.Node {
	var g = make([]graph.Node, 0)

	//q := `MATCH (n:User), (u:User) WHERE ID(n) = ` + strconv.Itoa(int(u.Id)) + ` AND NOT (n)-[]-(u)
	//RETURN u ORDER BY u.rating DESC SKIP ` + strconv.Itoa(skip) + ` LIMIT 50`

	return g
}

func getUser(Id int, Username string) (u User, err error) {

	var q string

	if Username != "" {
		q = `MATCH (u:User {username : "` + Username + `"}) RETURN u`
	} else {
		q = `MATCH (u:User) WHERE ID(u)= ` + strconv.Itoa(Id) + ` RETURN u`
	}

	data, _, _, _ := app.Neo.QueryNeoAll(q, nil)
	if len(data) == 0 {
		err = errors.New("Err : User doesn't exist")
		return
	} else {
		jso, _ := json.Marshal(data[0][0].(graph.Node).Properties)
		_ = json.Unmarshal(jso, &u)
		u.Id = data[0][0].(graph.Node).NodeIdentity
		return
	}
}
