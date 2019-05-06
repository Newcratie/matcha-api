package api

import (
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/errors"
	"strconv"
)

//MATCH (u:User), (n:User) WHERE ID(u) = 30 AND ID(n) = 238 return exists( (u)-[:LIKE]->(n) )
//MATCH (u:User) WHERE ID(u) = 30 MATCH (n:User) WHERE ID(n) = 238 CREATE (n)<-[:LIKE]-(u) return u, n
//MATCH (u)<-[r:LIKE]-(n) WHERE ID(u) = 30 AND ID(n) = 238 DELETE r
//MATCH (n)-[r:LIKE]-(u) WHERE ID(u) = 30 AND ID(n) = 238 DETACH DELETE r

func (app *App) dbMatchs(m Match) (valid bool, err error) {

	if app.dbExistBlocked(m) {
		err = errors.New("Blocked Relation")
		fmt.Println("****BLOCKED****")
		return false, err
	}

	if m.action == "LIKE" {
		app.dbCreateLike(m)
		fmt.Println("****LIKE****")
	} else if m.action == "DISLIKE" {
		app.dbCreateDislike(m)
		fmt.Println("****DISLIKE****")
	} else if m.action == "BLOCK" {
		app.dbCreateBlock(m)
		fmt.Println("****BLOCKED****")
	}
	return true, nil
}

func (app *App) dbExistBlocked(m Match) (valid bool) {

	prin("MAP == ", MapOf(m), "|")
	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` RETURN EXISTS( (u)-[:BLOCK]-(n) )`
	data, _, _, err := app.Neo.QueryNeoAll(ExistQuery, nil)
	if err != nil {
		prin("Err === ", err)
	}
	prin("DATA ===> ", data, "|")
	if data[0][0] == false {
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbCreateLike(m Match) (valid bool) {

	if app.dbExistRevLike(m) == false {
		MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` CREATE (u)-[:LIKE]->(n)`
		data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
		if err != nil {
			fmt.Println("*** CreateLike Query returned an Error ***", data)
			return false
		}
		return true
	} else {
		app.dbSetMatch(m)
	}
	return false
}

func (app *App) dbSetMatch(m Match) (valid bool) {

	app.dbDeleteRelation(m, "")
	MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` CREATE (u)-[:MATCH]->(n)`
	data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
	if err != nil {
		fmt.Println("*** Set MatchQuery returned an Error ***", data)
		return false
	}
	return true
}

func (app *App) dbCreateBlock(m Match) (valid bool) {

	app.dbDeleteRelation(m, "")
	MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` CREATE (u)-[:BLOCK]->(n)`
	prin("QUERY ==>", MatchQuery, "|")
	data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
	if err != nil {
		fmt.Println("*** CreateLike Query returned an Error ***", data)
		return false
	}
	return true
}

func (app *App) dbCreateDislike(m Match) (valid bool) {

	if app.dbExistRel(m, "DISLIKE") == false {
		MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` CREATE (u)-[:DISLIKE]->(n)`
		data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
		if err != nil {
			fmt.Println("*** CreateLike Query returned an Error ***", data)
			return false
		}
		return true
	}
	return false
}

func (app *App) dbExistRel(m Match, Rel string) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` RETURN EXISTS( (u)-[:` + Rel + `]->(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
	if data[0][0] == false {
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbExistDislike(m Match) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` RETURN EXISTS( (u)-[:DISLIKE]->(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
	if data[0][0] == false {
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbExistRevLike(m Match) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` RETURN EXISTS( (u)<-[:LIKE]-(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
	if data[0][0] == false {
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbDeleteRelation(m Match, Rel string) (valid bool) {

	if Rel != "" {
		DeleteQuery := `MATCH (n)-[t]-(u) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` DETACH DELETE t CREATE (u)<-[r:` + Rel + `']-(n)`
		data, _, _, err := app.Neo.QueryNeoAll(DeleteQuery, nil)
		if err != nil {
			fmt.Println("*** DeleteRelation Query returned an Error ***", data)
			return false
		}
	} else {
		DeleteQuery := `MATCH (u:User)-[r]-(n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` DETACH DELETE r`
		data, _, _, err := app.Neo.QueryNeoAll(DeleteQuery, nil)
		if err != nil {
			fmt.Println("*** DeleteRelation Query returned an Error ***", data)
			return false
		}
	}
	return true
}
