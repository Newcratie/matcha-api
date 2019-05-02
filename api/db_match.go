package api

import (
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/errors"
	"strconv"
)

//q := `MATCH (u:User) WHERE ID(u)={id} SET u.name = {username},
//u.username = {username}, u.password = {password},
//u.firstname = {firstname}, u.lastname = {lastname},
//u.birthday = {birthday}, u.random_token = {random_token},
//u.img1 = {img1}, u.img2 = {img2},
//u.img3 = {img3}, u.img4 = {img4},
//u.img5 = {img5}, u.biography = {biography},
//u.genre = {genre}, u.interest = {interest},
//u.city = {city}, u.zip = {zip},
//u.country = {country}, u.latitude = {latitude},
//u.longitude = {longitude}, u.geo_allowed = {geo_allowed},
//u.online = {online}, u.rating = {rating},
//u.email = {email}, u.access_lvl = {access_lvl})`

//MATCH (u:User), (n:User) WHERE ID(u) = 30 AND ID(n) = 238 return exists( (u)-[:LIKE]->(n) )
//MATCH (u:User) WHERE ID(u) = 30 MATCH (n:User) WHERE ID(n) = 238 CREATE (n)<-[:LIKE]-(u) return u, n
//MATCH (u)<-[r:LIKE]-(n) WHERE ID(u) = 30 AND ID(n) = 238 DELETE r
//MATCH (n)-[r:LIKE]-(u) WHERE ID(u) = 30 AND ID(n) = 238 DETACH DELETE r

func (app *App) dbMatchs(M Match) (valid bool, err error) {

	//fmt.Println("****IN DB MsssssssATCH****")
	if app.dbExistBlocked(M) {
		err = errors.New("Blocked Relation")
		return false, err
	}

	if M.Action == "LIKE" {
		app.dbCreateLike(M)
	} else if M.Action == "DISLIKE" {

	} else if M.Action == "BLOCK" {

	}

	//if Relation != "" {
	//	app.dbDeleteRelation(IdFrom, IdTo, Relation)
	//}
	//
	//if app.dbExistLike(IdFrom, IdTo, "LIKE") == true {
	//	if app.dbSetMatch(IdFrom, IdTo) == true {
	//		app.dbDeleteRelation(IdFrom, IdTo, "LIKE")
	//		return true
	//	}
	//} else if app.dbCreateLike(IdFrom, IdTo) == false {
	//	return false
	//}
	return true, nil
}

func (app *App) dbExistBlocked(M Match) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = {id_from} AND ID(n) = {id_to} RETURN EXISTS( (u)-[:BLOCKED]-(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, MapOf(M))
	if data[0][0] == false {
		//err = errors.New("wrong username or password")
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbCreateLike(M Match) (valid bool) {

	if app.dbExistLike(M) == false {
		MatchQuery := `MATCH (u:User), (n:User) WHERE ID(u) = {id_from} AND ID(n) = {id_to} CREATE (u)-[:LIKE]->(n)`
		data, _, _, err := app.Neo.QueryNeoAll(MatchQuery, MapOf(M))
		if err != nil {
			//err = errors.New("wrong username or password")
			fmt.Println("*** CreateLike Query returned an Error ***", data)
			return false
		}
		return true
	}
	return false
}

func (app *App) dbExistLike(M Match) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = {id_from} AND ID(n) = {id_to} RETURN EXISTS( (u)-[:LIKE]->(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, MapOf(M))
	if data[0][0] == false {
		//err = errors.New("wrong username or password")
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbExistDislike(M Match) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = {id_from} AND ID(n) = {id_to} RETURN EXISTS( (u)-[:DISLIKE]->(n) )`
	data, _, _, _ := app.Neo.QueryNeoAll(ExistQuery, MapOf(M))
	if data[0][0] == false {
		//err = errors.New("wrong username or password")
		fmt.Println("*** Exist Query returned FALSE ***")
		return false
	}
	return true
}

func (app *App) dbExistRevLike(IdFrom int, IdTo int, ExistRel string) (valid bool) {

	ExistQuery := `MATCH (u:User), (n:User) WHERE ID(u) = ` + strconv.Itoa(IdFrom) + ` AND ID(n) = ` + strconv.Itoa(IdTo) + ` RETURN EXISTS( (u)<-[:LIKE]-(n) )`
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
