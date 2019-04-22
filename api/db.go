package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

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

func (app *App) dbMatchs(IdFrom int, IdTo int) (valid bool) {

	if app.dbExistMatch(IdFrom, IdTo) == true {
		if app.dbSetMatch(IdFrom, IdTo) == false {
			return false
		}
	} else if app.dbCreateMatch(IdFrom, IdTo) == false {
		return false
	}
	return true
}

func (app *App) dbCreateMatch(IdFrom int, IdTo int) (valid bool) {

	MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(IdFrom) + ` AND ID(n) = ` + strconv.Itoa(IdTo) + ` CREATE (u)-[:LIKE]->(n)`
	data, _, _, _ := app.Neo.QueryNeoAll(MatchQuery, nil)
	if data[0][0] == false {
		//err = errors.New("wrong username or password")
		fmt.Println("*** MatchQuery returned false ***")
		return false
	}
	return true
}

func (app *App) dbExistMatch(IdFrom int, IdTo int) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(IdFrom) + ` AND ID(n) = ` + strconv.Itoa(IdTo) + ` RETURN EXISTS( (u)<-[:LIKE]-(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
	if data[0][0] == false {
		//err = errors.New("wrong username or password")
		fmt.Println("*** ExistQuery returned false ***")
		return false
	}
	return true
}

func (app *App) dbSetMatch(IdFrom int, IdTo int) (valid bool) {

	MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(IdFrom) + ` AND ID(n) = ` + strconv.Itoa(IdTo) + ` CREATE (u)-[:MATCH]->(n)`
	data, _, _, _ := app.Neo.QueryNeoAll(MatchQuery, nil)
	if data[0][0] == false {
		//err = errors.New("wrong username or password")
		fmt.Println("*** MatchQuery returned false ***")
		return false
	}
	return true
}

func (app *App) dbDeleteMatch(IdFrom int, IdTo int) (valid bool) {

	DeleteQuery := `MATCH (n)-[r:LIKE]-(u)  WHERE ID(u) = ` + strconv.Itoa(IdFrom) + ` AND ID(n) = ` + strconv.Itoa(IdTo) + ` DETACH DELETE r`
	data, _, _, _ := app.Neo.QueryNeoAll(DeleteQuery, nil)
	if data[0][0] == false {
		//err = errors.New("wrong username or password")
		fmt.Println("*** DeleteQuery returned false ***")
		return false
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

func (app *App) dbGetMatch(Id int, Filter *Filters) ([]graph.Node, error) {
	var g = make([]graph.Node, 0)
	var err error

	superQuery := `MATCH (u)-[m:MATCH]-(n) WHERE ID(u) = 30 return n`

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

//--------------------------------------------------------------------------------------------------------------//

func dbConnect() *sqlx.DB {
	connStr := "user=matcha password=secret dbname=matcha host=" +
		os.Getenv("POSTGRES_HOST") +
		" port=5432 sslmode=disable"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

//func (app *App) ValidToken(randomToken string) error {
//	var u User
//	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:Person) WHERE n.random_token = "`+ randomToken + `" RETURN n`, nil)
//	if len(data) == 0 {
//		return errors.New("Invalid Link")
//	} else if u.AccessLvl == 1 {
//		return errors.New("Email already validated")
//	}
//	//_, err = app.Db.NamedExec(`UPDATE "public"."users" SET "access_lvl" = 1 WHERE "id" = :id`, u)
//	return nil
//}

func MapOf(u interface{}) (m map[string]interface{}) {
	m = make(map[string]interface{})
	jso, _ := json.Marshal(u)
	_ = json.Unmarshal(jso, &m)
	return m
}
