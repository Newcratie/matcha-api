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

func (app *App) dbGetPeople(Id int, Filter *Filters) ([]graph.Node, error) {

	fmt.Println("voilqaa ===>")
	fmt.Println(Filter)
	fmt.Println("voila2 ===>")
	fmt.Printf("%+v\n", Filter)

	var g []graph.Node
	var err error
	//tempQuery := `MATCH (n:User) WHERE ID(n) <> ` + strconv.Itoa(Id) + ` RETURN n LIMIT 40`
	superQuery := customQuery(Id, Filter)
	//`MATCH (n:User) WHERE ID(n) <> ` + strconv.Itoa(Id) + ` RETURN n LIMIT 40`
	// for age == MATCH (u:User) WHERE u.birthday > "1914-10-06T18:51:39.178248882Z" AND u.birthday < "1916-10-06T18:51:39.178248882Z" return u
	//  for score/rating == MATCH (u:User) WHERE u.rating > 0 AND u.rating < 5 return u
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

//MATCH (u) WHERE (u.latitude > ` + strconv.Itoa(Filter.Location[0]) + ` AND u.longitude < ` + strconv.Itoa(Filter.Location[1]) + `)
//MATCH (u) WHERE (u.birthday > ` + maxAge + ` AND u.birthday < ` + minAge + `)

func customQuery(Id int, Filter *Filters) (superQuery string) {

	minAge := ageConvert(Id, Filter.Age[0])
	maxAge := ageConvert(Id, Filter.Age[1])

	superQuery = `MATCH (u:User) WHERE (u.rating > ` + strconv.Itoa(Filter.Score[0]) + ` AND u.rating < ` + strconv.Itoa(Filter.Score[1]) + `)
	MATCH (u) WHERE (u.birthday > "` + maxAge + `" AND u.birthday < "` + minAge + `")
	return u`

	return superQuery
}

func ageConvert(Id int, Age int) (birthYear string) {

	//p := fmt.Println

	//t := time.Now()
	//p("With Format", t.Format(time.RFC3339))

	fmt.Println(Age)
	now := time.Now()
	fmt.Println("Before Format : ", now)
	now = now.AddDate(-(Age), 0, 0)
	//now, _ = time.Parse(timeFormat, "1906-12-27T17:14:59.681469185Z")
	fmt.Println("After Format : ", now)
	birthYear = now.Format(time.RFC3339Nano)
	//birthYear = strings.Replace(birthYear, " +0000 UTC", "", -1)
	fmt.Println("On STRINGED : ", birthYear)
	//fmt.Println("user ID : ", Id)
	//userDateOfBirth, _, _, _ := app.Neo.QueryNeoAll(`MATCH (u:User) WHERE Id(u)=` + strconv.Itoa(67) + ` return u.birthday`, nil)
	//fmt.Println("userbirth : ", userDateOfBirth)
	//now := time.Now()
	//fmt.Println("Time now : ", now)

	return birthYear
}

//func kmConvert(Km int) (lat1 string, lon1 string) {
//
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
