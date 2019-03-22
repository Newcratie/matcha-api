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
)

func (app *App) insertUser(u User) {
	q := `CREATE (u:User{name: {username},
username:{username}, password:{password},
firstname:{firstname}, lastname:{lastname},
birthday:{birthday}, random_token: {random_token},
email: {email}})`
	st := app.prepareStatement(q)
	executeStatement(st, MapOf(u))
}

func (app *App) getUser(Username string) (u User, err error) {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User{username : "`+Username+`"}) SET n.online = true RETURN id(n), n`, nil)
	if len(data) == 0 {
		err = errors.New("wrong username or password")
		return
	} else {
		jso, _ := json.Marshal(data[0][0].(graph.Node).Properties)
		_ = json.Unmarshal(jso, &u)
		return
	}
}

func (app *App) getBasicUser(Id int) (u User, err error) {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User) WHERE id(n)= `+string(Id)+`}) RETURN id(n), n`, nil)
	fmt.Println("basic: ", data)
	if len(data) == 0 {
		err = errors.New("wrong username or password")
		return
	} else {
		jso, _ := json.Marshal(data[0][0].(graph.Node).Properties)
		_ = json.Unmarshal(jso, &u)
		return
	}
}

func (app *App) getBasicDates(Id int) (u []graph.Node, err error) {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User) RETURN n`, nil)
	fmt.Println("dates: ", data)
	if len(data) == 0 {
		err = errors.New("wrong username or password")
		return
	} else {
		jso, _ := json.Marshal(data[0])
		_ = json.Unmarshal(jso, &u)
		fmt.Println("graph.Node ===> ", u)
		return
	}
}
func (app *App) accountExist(rf registerForm) bool {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:Person) WHERE n.username = "`+rf.Username+`" RETURN n`, nil)
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
