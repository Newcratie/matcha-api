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

func (app *App) dbMatchs(M Match) (valid bool, err error) {

	//fmt.Println("****IN DB MsssssssATCH****")
	if app.dbExistBlocked(M) {
		err = errors.New("Blocked Relation")
		fmt.Println("****BLOCKED****")
		return false, err
	}

	if M.Action == "LIKE" {
		app.dbCreateLike(M)
		fmt.Println("****LIKE****")
	} else if M.Action == "DISLIKE" {
		app.dbCreateDislike(M)
		fmt.Println("****DISLIKE****")
	} else if M.Action == "BLOCKED" {
		app.dbCreateBlock(M)
		fmt.Println("****BLOCKED****")
	}
	return true, nil
}

func (app *App) dbExistBlocked(M Match) (valid bool) {

	prin("MAP == ", MapOf(M), "|")
	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(M.IdFrom) + ` AND ID(n) = ` + strconv.Itoa(M.IdTo) + ` RETURN EXISTS( (u)-[:BLOCKED]-(n) )`
	data, _, _, err := app.Neo.QueryNeoAll(ExistQuery, nil)
	if err != nil {
		prin("Err === ", err)
	}
	prin("DATA ===> ", data, "|")
	if data[0][0] == false {
		//err = errors.New("wrong username or password")
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbCreateLike(M Match) (valid bool) {

	if app.dbExistRevLike(M) == false {
		MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(M.IdFrom) + ` AND ID(n) = ` + strconv.Itoa(M.IdTo) + ` CREATE (u)-[:LIKE]->(n)`
		data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
		if err != nil {
			//err = errors.New("wrong username or password")
			fmt.Println("*** CreateLike Query returned an Error ***", data)
			return false
		}
		return true
	} else {
		app.dbSetMatch(M)
	}
	return false
}

func (app *App) dbSetMatch(M Match) (valid bool) {

	app.dbDeleteRelation(M, "")
	MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(M.IdFrom) + ` AND ID(n) = ` + strconv.Itoa(M.IdTo) + ` CREATE (u)-[:MATCH]->(n)`
	data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
	if err != nil {
		fmt.Println("*** Set MatchQuery returned an Error ***", data)
		return false
	}
	return true
}

func (app *App) dbCreateBlock(M Match) (valid bool) {

	app.dbDeleteRelation(M, "")
	MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(M.IdFrom) + ` AND ID(n) = ` + strconv.Itoa(M.IdTo) + ` CREATE (u)-[:BLOCKED]->(n)`
	prin("QUERY ==>", MatchQuery, "|")
	data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
	if err != nil {
		fmt.Println("*** CreateLike Query returned an Error ***", data)
		return false
	}
	return true
}

func (app *App) dbCreateDislike(M Match) (valid bool) {

	if app.dbExistRel(M, "DISLIKE") == false {
		MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(M.IdFrom) + ` AND ID(n) = ` + strconv.Itoa(M.IdTo) + ` CREATE (u)-[:DISLIKE]->(n)`
		data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
		if err != nil {
			fmt.Println("*** CreateLike Query returned an Error ***", data)
			return false
		}
		return true
	}
	return false
}

func (app *App) dbExistRel(M Match, Rel string) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(M.IdFrom) + ` AND ID(n) = ` + strconv.Itoa(M.IdTo) + ` RETURN EXISTS( (u)-[:` + Rel + `]->(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
	if data[0][0] == false {
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbExistDislike(M Match) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(M.IdFrom) + ` AND ID(n) = ` + strconv.Itoa(M.IdTo) + ` RETURN EXISTS( (u)-[:DISLIKE]->(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
	if data[0][0] == false {
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbExistRevLike(M Match) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(M.IdFrom) + ` AND ID(n) = ` + strconv.Itoa(M.IdTo) + ` RETURN EXISTS( (u)<-[:LIKE]-(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
	if data[0][0] == false {
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbDeleteRelation(M Match, Rel string) (valid bool) {

	if Rel != "" {
		DeleteQuery := `MATCH (n)-[t]-(u) WHERE ID(u) = ` + strconv.Itoa(M.IdFrom) + ` AND ID(n) = ` + strconv.Itoa(M.IdTo) + ` DETACH DELETE t CREATE (u)<-[r:` + Rel + `']-(n)`
		data, _, _, err := app.Neo.QueryNeoAll(DeleteQuery, nil)
		if err != nil {
			fmt.Println("*** DeleteRelation Query returned an Error ***", data)
			return false
		}
	} else {
		DeleteQuery := `MATCH (u:User)-[r]-(n:User) WHERE ID(u) = ` + strconv.Itoa(M.IdFrom) + ` AND ID(n) = ` + strconv.Itoa(M.IdTo) + ` DETACH DELETE r`
		data, _, _, err := app.Neo.QueryNeoAll(DeleteQuery, nil)
		if err != nil {
			fmt.Println("*** DeleteRelation Query returned an Error ***", data)
			return false
		}
	}
	return true
}
