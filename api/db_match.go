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

	if m.action == "LIKE" && app.dbExistRel(m, like) == false {
		app.dbCreateLike(m)
		fmt.Println("****LIKE****")
	} else if m.action == "DISLIKE" && app.dbExistRel(m, dislike) == false {
		app.dbCreateDislike(m)
		fmt.Println("****DISLIKE****")
	} else if m.action == "BLOCK" {
		app.dbCreateBlock(m)
		fmt.Println("****BLOCKED****")
	}
	return true, nil
}

func (app *App) dbExistBlocked(m Match) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` RETURN EXISTS( (u)-[:BLOCK]-(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
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
	MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` CREATE (u)-[:MATCHED]->(n)`
	data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
	if err != nil {
		fmt.Println("*** Set MatchQuery returned an Error ***", data)
		return false
	}
	return true
}

func (app *App) dbCreateBlock(m Match) (valid bool) {

	app.dbDeleteDirectionalRelation(m, "")
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

	if app.dbExistMatch(m) == true {
		app.dbDeleteRelation(m, matched)
	}
	app.dbDeleteDirectionalRelation(m, "")
	MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` CREATE (u)-[:DISLIKE]->(n)`
	data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, nil)
	if err != nil {
		fmt.Println("*** CreateLike Query returned an Error ***", data)
		return false
	}
	return true
}

func (app *App) dbExistMatch(m Match) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` RETURN EXISTS( (u)-[:MATCHED]-(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
	if data[0][0] == false {
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
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

func (app *App) dbExistRevLike(m Match) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` RETURN EXISTS( (u)<-[:LIKE]-(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, nil)
	if data[0][0] == false {
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbDeleteDirectionalRelation(m Match, Rel string) (valid bool) {

	if Rel != "" {
		DeleteQuery := `MATCH (u)-[t:` + Rel + `]->(n) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` DETACH DELETE t`
		data, _, _, err := app.Neo.QueryNeoAll(DeleteQuery, nil)
		if err != nil {
			fmt.Println("*** DeleteRelation Query returned an Error ***", data)
			return false
		}
	} else {
		DeleteQuery := `MATCH (u:User)-[r]->(n:User) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` DETACH DELETE r`
		data, _, _, err := app.Neo.QueryNeoAll(DeleteQuery, nil)
		if err != nil {
			fmt.Println("*** DeleteRelation Query returned an Error ***", data)
			return false
		}
	}
	return true
}

func (app *App) dbDeleteRelation(m Match, Rel string) (valid bool) {

	if Rel != "" {
		DeleteQuery := `MATCH (n)-[r:` + Rel + `]-(u) WHERE ID(u) = ` + strconv.Itoa(m.idFrom) + ` AND ID(n) = ` + strconv.Itoa(m.idTo) + ` DETACH DELETE r`
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
