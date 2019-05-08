package api

import (
	"encoding/json"
	"errors"
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	_ "github.com/lib/pq"
	"strconv"
	"time"
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
	To        int64  `json:"to"`
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
			int64(suitorId),
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
city:{city}, zip: {zip},
country:{country}, latitude: {latitude},
longitude:{longitude}, geo_allowed: {geo_allowed},
online:{online}, rating: {rating},
email: {email}, access_lvl: 1, last_conn: {last_conn},
ilike: {ilike}, relation: {relation}})`
	//fmt.Println("Query == ", q)
	st := app.prepareStatement(q)
	executeStatement(st, MapOf(u))
	return
}

func (app *App) updateUser(u User) {
	//fmt.Println(MapOf(u))
	var id string
	id = strconv.Itoa(int(u.Id))
	q := `MATCH (u:User) WHERE ID(u) = ` + id + ` SET u.name = {username},
	u.firstname = {firstname}, u.lastname = {lastname},
	u.username = {username}, u.password = {password},
	u.birthday = {birthday}, u.random_token = {random_token},
	u.img1 = {img1}, u.img2 = {img2},
	u.img3 = {img3}, u.img4 = {img4},
	u.img5 = {img5}, u.biography = {biography},
	u.genre = {genre}, u.interest = {interest},
	u.city = {city}, u.zip = {zip},
	u.country = {country}, u.latitude = {latitude},
	u.longitude = {longitude}, u.geo_allowed = {geo_allowed},
	u.online = {online}, u.rating = {rating},
	u.email = {email}, u.access_lvl = {access_lvl},
	u.tags = {tags},  u.last_conn = {last_conn}`
	st := app.prepareStatement(q)
	executeStatement(st, MapOf(u))
	return
}

func (app *App) updateLastConn(u User) {
	prin("**** IN UPDATE CONN ****")
	id := strconv.Itoa(int(u.Id))
	q := `MATCH (u:User) WHERE ID(u) = ` + id + ` SET u.online = {online}, u.last_conn = {last_conn}`
	st := app.prepareStatement(q)
	executeStatement(st, MapOf(u))
	return
}

func (app *App) getUser(Id int, Username string) (u User, err error) {

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

	superQuery := `MATCH (u)-[m:MATCHED]-(n) WHERE ID(u) = ` + strconv.Itoa(Id) + ` return n`

	data, _, _, _ := app.Neo.QueryNeoAll(superQuery, nil)

	if len(data) == 0 {
		err = errors.New("Err : No Match found for given User ID")
		return g, err
	} else {
		for _, d := range data {
			g = append(g, d[0].(graph.Node))
		}
		return g, err

	}
}

func (app *App) dbGetUserProfile(Id int) (graph.Node, error) {
	var g = graph.Node{}
	var err error

	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User) WHERE ID(n) = `+strconv.Itoa(Id)+` SET n.online = true RETURN  n`, nil)
	if len(data) == 0 {
		err = errors.New("Err : User Id doesn't exist")
		return g, err
	} else {
		g = data[0][0].(graph.Node)
		delete(g.Properties, "password")
		return g, err
	}
}

func (app *App) dbGetPeople(Id int, Filter *Filters) ([]graph.Node, error) {
	var g = make([]graph.Node, 0)
	var err error

	var m Match
	m.action = like
	// A custom query with applied Filters
	time.Sleep(500 * time.Millisecond)
	superQuery := customQuery(Id, Filter)

	data, _, _, err := app.Neo.QueryNeoAll(superQuery, nil)
	u, _ := app.getUser(Id, "")
	if len(data) == 0 {
		err = errors.New("err : filters doesn't match anyone")
		return g, err
	} else {
		for _, d := range data {
			m.idFrom = int(d[0].(graph.Node).NodeIdentity)
			m.idTo = Id
			d[0].(graph.Node).Properties["ilike"] = app.dbExistRel(m, m.action)
			d[0].(graph.Node).Properties["relation"] = app.dbGetRelationType(m)
			lonTo, _ := getFloat(d[0].(graph.Node).Properties["longitude"])
			latTo, _ := getFloat(d[0].(graph.Node).Properties["latitude"])
			Genre, _ := d[0].(graph.Node).Properties["genre"].(string)
			Interest, _ := d[0].(graph.Node).Properties["interest"].(string)
			delete(d[0].(graph.Node).Properties, "password")

			// Haversine will return the distance between 2 Lat/Lon in Kilometers

			if Haversine(u.Longitude, u.Latitude, lonTo, latTo) <= Filter.Location[1] {
				if valid := setInterest(Genre, Interest, Id); valid == true {
					g = append(g, d[0].(graph.Node))
				}
			}
		}
		return g, err
	}
}

func (app *App) usernameExist(Username string) bool {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User {username: {username}}) RETURN n`, map[string]interface{}{"username": Username})
	if len(data) == 0 {
		return false
	} else {
		return true
	}
}

func (app *App) emailExist(Email string) bool {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User {email: {email}}) RETURN n`, map[string]interface{}{"email": Email})
	if len(data) == 0 {
		return false
	} else {
		return true
	}
}

func (app *App) prepareStatement(query string) bolt.Stmt {
	conn, err := app.Db.OpenPool()
	st, err := conn.PrepareNeo(query)
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
