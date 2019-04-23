package api

import (
	"encoding/json"
	"errors"
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	_ "github.com/lib/pq"
	"strconv"
)

func (app *App) insertMessage(byt []byte) {
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	dat["author"] = int64(dat["author"].(float64))
	dat["to"] = int64(dat["to"].(float64))
	dat["timestamp"] = int64(dat["timestamp"].(float64))
	q := `
MATCH (a:User),(b:User)
WHERE ID(a)={author} AND ID(b)={to}
CREATE (a)-[s:SAYS]->(message:Message {msg:{msg}, author: {author}, id:{id}, timestamp:{timestamp}})-[t:TO]->(b)`
	st := app.prepareStatement(q)
	executeStatement(st, dat)
}

type Messages struct {
	Id        int64  `json:"id"`
	Author    int64  `json:"author"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (app *App) dbGetMessages(userId, suitorId int) ([]Messages, error) {
	q := `
MATCH (a:User)-[]-(n:Message)-[]-(b:User) 
WHERE ID(a)={user_id} AND ID(b)={suitor_id}
RETURN n
ORDER BY ID(n)
`
	msgs := make([]Messages, 0)

	data, _, _, _ := app.Neo.QueryNeoAll(q, map[string]interface{}{"user_id": userId, "suitor_id": suitorId})
	fmt.Println(data)
	for _, tab := range data {
		msgs = append(msgs, Messages{
			int64(tab[0].(graph.Node).Properties["id"].(float64)),
			tab[0].(graph.Node).Properties["author"].(int64),
			tab[0].(graph.Node).Properties["msg"].(string),
			tab[0].(graph.Node).Properties["timestamp"].(int64),
		})
	}
	return msgs, nil
}

func (app *App) insertUser(u User) {
	fmt.Println(MapOf(u))
	q := `CREATE (u:User{name: {username},
username:{username}, password:{password},
firstname:{firstname}, lastname:{lastname},
birthday:{birthday}, random_token: {random_token},
img1:{img1}, img2: {img2},
img3:{img3}, img4: {img4},
img5:{img5}, biography: {biography},
genre:{genre}, interest: {interest},
img5:{img5}, biography: {biography},
city:{city}, zip: {zip},
country:{country}, latitude: {latitude},
longitude:{longitude}, geo_allowed: {geo_allowed},
online:{online}, rating: {rating},
email: {email}, access_lvl: 0})`
	st := app.prepareStatement(q)
	executeStatement(st, MapOf(u))
}

//MATCH (u:User), (n:User) WHERE ID(u) = 30 AND ID(n) = 238 return exists( (u)-[:LIKE]->(n) )
//MATCH (u:User) WHERE ID(u) = 30 MATCH (n:User) WHERE ID(n) = 238 CREATE (n)<-[:LIKE]-(u) return u, n
//MATCH (u)<-[r:LIKE]-(n) WHERE ID(u) = 30 AND ID(n) = 238 DELETE r
//MATCH (n)-[r:LIKE]-(u) WHERE ID(u) = 30 AND ID(n) = 238 DETACH DELETE r

func (app *App) dbMatchs(IdFrom int, IdTo int, Relation string) (valid bool) {

	//fmt.Println("****IN DB MsssssssATCH****")

	if Relation != "" {
		app.dbDeleteRelation(IdFrom, IdTo, Relation)
	}

	if app.dbExistLike(IdFrom, IdTo, "LIKE") == true {
		if app.dbSetMatch(IdFrom, IdTo) == true {
			app.dbDeleteRelation(IdFrom, IdTo, "LIKE")
			return true
		}
	} else if app.dbCreateLike(IdFrom, IdTo) == false {
		return false
	}
	return false
}

func (app *App) dbCreateLike(IdFrom int, IdTo int) (valid bool) {

	if app.dbExistLike(IdFrom, IdTo, "DISLIKE") == false {
		MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(IdFrom) + ` AND ID(n) = ` + strconv.Itoa(IdTo) + ` CREATE (u)-[:LIKE]->(n)`
		data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
		if err != nil {
			//err = errors.New("wrong username or password")
			fmt.Println("*** CreateLike Query returned an Error ***", data)
			return false
		}
	}
	return true
}

func (app *App) dbExistLike(IdFrom int, IdTo int, ExistRel string) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(IdFrom) + ` AND ID(n) = ` + strconv.Itoa(IdTo) + ` RETURN EXISTS( (u)<-[:` + ExistRel + `]-(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
	if data[0][0] == false {
		//err = errors.New("wrong username or password")
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbSetMatch(IdFrom int, IdTo int) (valid bool) {

	MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(IdFrom) + ` AND ID(n) = ` + strconv.Itoa(IdTo) + ` CREATE (u)-[:MATCH]->(n)`
	data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
	if err != nil {
		//err = errors.New("wrong username or password")
		fmt.Println("*** Set MatchQuery returned an Error ***", data)
		return false
	}
	return true
}

func (app *App) dbDeleteRelation(IdFrom int, IdTo int, Rel string) (valid bool) {

	if Rel == "DISLIKE" {
		DeleteQuery := `MATCH (n)-[t]-(u) WHERE ID(u) = ` + strconv.Itoa(IdFrom) + ` AND ID(n) = ` + strconv.Itoa(IdTo) + ` DETACH DELETE t CREATE (u)<-[r:DISLIKE]-(n)`
		data, _, _, err := app.Neo.QueryNeoAll(DeleteQuery, nil)
		if err != nil {
			//err = errors.New("wrong username or password")
			fmt.Println("*** DeleteRelation Query returned an Error ***", data)
			return false
		}
	} else {
		DeleteQuery := `MATCH (n)-[m:` + Rel + `]-(u)  WHERE ID(u) = ` + strconv.Itoa(IdFrom) + ` AND ID(n) = ` + strconv.Itoa(IdTo) + ` DETACH DELETE m`
		data, _, _, err := app.Neo.QueryNeoAll(DeleteQuery, nil)
		if err != nil {
			//err = errors.New("wrong username or password")
			fmt.Println("*** DeleteRelation Query returned an Error ***", data)
			return false
		}
	}
	return true
}

func (app *App) getUser(Username string) (u User, err error) {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User{username : "`+Username+`"}) SET n.online = true RETURN  n`, nil)
	fmt.Println(data)
	if len(data) == 0 {
		err = errors.New("wrong username or password")
		return
	} else {
		jso, _ := json.Marshal(data[0][0].(graph.Node).Properties)
		_ = json.Unmarshal(jso, &u)
		u.Id = data[0][0].(graph.Node).NodeIdentity
		return
	}
}

func (app *App) getBasicUser(Id int) (u User, err error) {
	data, _, _, err := app.Neo.QueryNeoAll(`MATCH (n:User) WHERE id(n) = `+strconv.Itoa(Id)+` RETURN n`, nil)
	fmt.Println("basic: ", data)
	if len(data) == 0 || err != nil {
		return
	} else {
		jso, _ := json.Marshal(data[0][0].(graph.Node))
		_ = json.Unmarshal(jso, &u)
		return
	}
}

func (app *App) dbGetMatchs(Id int) ([]graph.Node, error) {
	var g = make([]graph.Node, 0)
	var err error

	superQuery := `MATCH (u)-[m:MATCH]-(n) WHERE ID(u) = ` + strconv.Itoa(Id) + ` return n`

	data, _, _, _ := app.Neo.QueryNeoAll(superQuery, nil)

	if len(data) == 0 {
		err = errors.New("wrong username or password")
		return g, err
	} else {
		for _, d := range data {
			g = append(g, d[0].(graph.Node))
		}
		fmt.Println("YOOOOOOOO ===")
		fmt.Println(g)
		return g, err

	}
}

func (app *App) dbGetUserProfile(Id int) ([]graph.Node, error) {
	var g = make([]graph.Node, 0)
	var err error

	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User) WHERE ID(n) = `+strconv.Itoa(Id)+` SET n.online = true RETURN  n`, nil)
	fmt.Println(data)
	if len(data) == 0 {
		err = errors.New("wrong username or password")
		return g, err
	} else {
		for _, d := range data {
			g = append(g, d[0].(graph.Node))
		}
		fmt.Println("USER Info = ", g)
		return g, err
	}
}

func (app *App) dbGetPeople(Id int, Filter *Filters) ([]graph.Node, error) {
	var g = make([]graph.Node, 0)
	var err error

	// A custom query with applied Filters
	superQuery := customQuery(Id, Filter)

	data, _, _, _ := app.Neo.QueryNeoAll(superQuery, nil)

	if len(data) == 0 {
		err = errors.New("wrong username or password")
		return g, err
	} else {
		for _, d := range data {
			lonTo, _ := getFloat(d[0].(graph.Node).Properties["longitude"])
			latTo, _ := getFloat(d[0].(graph.Node).Properties["latitude"])

			// Haversine will return the distance between 2 Lat/Lon in Kilometers

			if Haversine(0, 0, lonTo, latTo) <= Filter.Location[1] {
				g = append(g, d[0].(graph.Node))
			}
		}
		return g, err
	}
}

func (app *App) usernameExist(rf registerForm) bool {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User {username: {username}}) RETURN n`, map[string]interface{}{"username": rf.Username})
	if len(data) == 0 {
		return false
	} else {
		return true
	}
}

func (app *App) emailExist(rf registerForm) bool {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User {email: {email}}) RETURN n`, map[string]interface{}{"email": rf.Email})
	if len(data) == 0 {
		return false
	} else {
		return true
	}
}

func (app *App) prepareStatement(query string) bolt.Stmt {
	st, err := app.Neo.PrepareNeo(query)
	handleError(err)
	return st
}

func executeStatement(st bolt.Stmt, m map[string]interface{}) {
	result, err := st.ExecNeo(m)
	handleError(err)
	_, err = result.RowsAffected()
	handleError(err)

	st.Close()
}
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func MapOf(u interface{}) (m map[string]interface{}) {
	m = make(map[string]interface{})
	jso, _ := json.Marshal(u)
	_ = json.Unmarshal(jso, &m)
	return m
}
