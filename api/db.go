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
	"time"
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
	var g []graph.Node
	var err error

	// A custom query with applied Filters
	superQuery := `MATCH (u:User) RETURN u LIMIT 10`

	data, _, _, _ := app.Neo.QueryNeoAll(superQuery, nil)

	if len(data) == 0 {
		err = errors.New("wrong username or password")
		return g, err
	} else {
		for _, d := range data {
			g = append(g, d[0].(graph.Node))
		}
		return g, err
	}

}

func (app *App) dbGetPeople(Id int, Filter *Filters) ([]graph.Node, error) {
	var g []graph.Node
	var err error

	// A custom query with applied Filters
	superQuery := customQuery(Id, Filter)

	data, _, _, _ := app.Neo.QueryNeoAll(superQuery, nil)

	if len(data) == 0 {
		err = errors.New("wrong username or password")
		return g, err
	} else {
		for _, d := range data {
			g = append(g, d[0].(graph.Node))
		}
		return g, err
	}
}

// Miss Latidude/Longitude Max/Min
func customQuery(Id int, Filter *Filters) (superQuery string) {

	minAge := ageConvert(Filter.Age[0])
	maxAge := ageConvert(Filter.Age[1])
	minLat, maxLat, minLon, maxLon := latLongCheck(Id, Filter.Location[1])

	fmt.Println("Salut c'est COOL : ", minLat, maxLat, minLon, maxLon)

	superQuery = `MATCH (u:User) WHERE (u.rating > ` + strconv.Itoa(Filter.Score[0]) + ` AND u.rating < ` + strconv.Itoa(Filter.Score[1]) + `)
	MATCH (u) WHERE (u.birthday > "` + maxAge + `" AND u.birthday < "` + minAge + `")
	return u`

	return superQuery
}

func ageConvert(Age int) (birthYear string) {

	now := time.Now()
	now = now.AddDate(-(Age), 0, 0)
	birthYear = now.Format(time.RFC3339Nano)
	return birthYear
}

func latLongCheck(Id int, Km int) (lat1 string, lon1 string, lat2 string, lon2 string) {
	Id = 42
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (s) WHERE ID(s)=`+strconv.Itoa(Id)+` return s.latitude, s.longitude`, nil)

	fmt.Println("DATA : ", data, Id)

	return "0", "0", "0", "0"
}

//func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
//	const PI float64 = 3.141592653589793
//
//	radlat1 := float64(PI * lat1 / 180)
//	radlat2 := float64(PI * lat2 / 180)
//
//	theta := float64(lng1 - lng2)
//	radtheta := float64(PI * theta / 180)
//
//	dist := math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta)
//
//	if dist > 1 {
//		dist = 1
//	}
//
//	dist = math.Acos(dist)
//	dist = dist * 180 / PI
//	dist = dist * 60 * 1.1515
//
//	if len(unit) > 0 {
//		if unit[0] == "K" {
//			dist = dist * 1.609344
//		} else if unit[0] == "N" {
//			dist = dist * 0.8684
//		}
//	}
//
//	return dist
//}

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
