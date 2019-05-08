package main

import (
	"github.com/Newcratie/matcha-api/api"
	"github.com/Newcratie/matcha-api/api/hash"
	"github.com/brianvoe/gofakeit"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"strconv"
	"strings"
	"time"
)

type App api.App

var app *App

func newRandomMale() api.User {
	var f *gofakeit.PersonInfo
	max := 1
	f = gofakeit.Person()
	interest := make([]string, 3)
	interest[0] = "bi"
	interest[1] = "hetero"
	interest[2] = "homo"
	tagtab := make([]string, 1)
	//tagg := []Tag
	for i := 0; i < max; i++ {
		tagtab[i] = gofakeit.Color()
		//tagg.
	}

	Latitude, _ := gofakeit.LatitudeInRange(42.490627, 50.264989)
	Longitude, _ := gofakeit.LongitudeInRange(-3.396493, 9.517944)
	return api.User{Username: gofakeit.Username(),
		Password:  hash.Encrypt(api.HashKey, "'"),
		FirstName: f.FirstName,
		LastName:  f.LastName,
		Email:     gofakeit.Email(),
		Img1:      "https://randomuser.me/api/portraits/men/" + strconv.Itoa(gofakeit.Number(1, 45)) + ".jpg",
		Img2:      "https://randomuser.me/api/portraits/men/" + strconv.Itoa(gofakeit.Number(1, 45)) + ".jpg",
		Img3:      gofakeit.ImageURL(300, 300),
		Img4:      gofakeit.ImageURL(300, 300),
		Img5:      gofakeit.ImageURL(300, 300),
		Biography: gofakeit.Paragraph(1, 4, 10, " "),
		Birthday: gofakeit.DateRange(time.Date(1900, 01, 01, 00, 00, 00, 00, time.UTC),
			time.Date(2000, 01, 01, 00, 00, 00, 00, time.UTC)),
		Genre:      "male",
		Interest:   gofakeit.RandString(interest),
		AccessLvl:  1,
		Online:     gofakeit.Bool(),
		Rating:     gofakeit.Number(0, 100),
		City:       gofakeit.City(),
		Zip:        gofakeit.Zip(),
		Country:    gofakeit.Country(),
		Latitude:   Latitude,
		Longitude:  Longitude,
		GeoAllowed: gofakeit.Bool(),
		CreatedAt: gofakeit.DateRange(time.Date(1900, 01, 01, 00, 00, 00, 00, time.Local),
			time.Date(2017, 01, 01, 00, 00, 00, 00, time.Local)),
		Tags: tagtab,
		LastConn: gofakeit.DateRange(time.Date(2016, 01, 01, 00, 00, 00, 00, time.Local),
			time.Date(2017, 01, 01, 00, 00, 00, 00, time.Local)),
		Ilike:    false,
		Relation: "none",
	}
}

func newRandomFemale() api.User {
	Latitude, _ := gofakeit.LatitudeInRange(42.490627, 50.264989)
	Longitude, _ := gofakeit.LongitudeInRange(-3.396493, 9.517944)
	var f *gofakeit.PersonInfo
	max := 1
	f = gofakeit.Person()
	interest := make([]string, 3)
	interest[0] = "bi"
	interest[1] = "hetero"
	interest[2] = "homo"
	tagtab := make([]string, 1)
	for i := 0; i < max; i++ {
		tagtab[i] = gofakeit.Color()
	}
	return api.User{Username: gofakeit.Username(),
		Password:  hash.Encrypt(api.HashKey, "'"),
		FirstName: f.FirstName,
		LastName:  f.LastName,
		Email:     gofakeit.Email(),
		Img1:      "https://randomuser.me/api/portraits/women/" + strconv.Itoa(gofakeit.Number(1, 45)) + ".jpg",
		Img2:      "https://randomuser.me/api/portraits/women/" + strconv.Itoa(gofakeit.Number(1, 45)) + ".jpg",
		Img3:      gofakeit.ImageURL(300, 300),
		Img4:      gofakeit.ImageURL(300, 300),
		Img5:      gofakeit.ImageURL(300, 300),
		Biography: gofakeit.Paragraph(1, 4, 10, " "),
		Birthday: gofakeit.DateRange(time.Date(1900, 01, 01, 00, 00, 00, 00, time.UTC),
			time.Date(2000, 01, 01, 00, 00, 00, 00, time.UTC)),
		Genre:      "female",
		Interest:   gofakeit.RandString(interest),
		AccessLvl:  1,
		Online:     gofakeit.Bool(),
		Rating:     gofakeit.Number(0, 100),
		City:       gofakeit.City(),
		Zip:        gofakeit.Zip(),
		Country:    gofakeit.Country(),
		Latitude:   Latitude,
		Longitude:  Longitude,
		GeoAllowed: gofakeit.Bool(),
		CreatedAt: gofakeit.DateRange(time.Date(1900, 01, 01, 00, 00, 00, 00, time.Local),
			time.Date(2017, 01, 01, 00, 00, 00, 00, time.Local)),
		Tags: tagtab,
		LastConn: gofakeit.DateRange(time.Date(2016, 01, 01, 00, 00, 00, 00, time.Local),
			time.Date(2017, 01, 01, 00, 00, 00, 00, time.Local)),
		Ilike:    false,
		Relation: "none",
	}
}

func main() {
	const max = 260

	app.Db, _ = bolt.NewDriverPool("bolt://neo4j:secret@neo4j:7687", 1000)
	app.Neo, _ = app.Db.OpenPool()
	defer app.Neo.Close()
	for i := 0; i < max; i++ {
		s := gofakeit.Color()
		s = strings.ToLower(s)
		app.Neo.QueryNeoAll(`MERGE (t:TAG {key: "`+s+`", text: "#`+strings.Title(s)+`", value: "`+s+`"}) `, nil)
	}
	for i := 0; i < max; i++ {
		u := newRandomMale()
		app.insertFakeUser(u)
		AddTagRelation(u)
		u = newRandomFemale()
		app.insertFakeUser(u)
		AddTagRelation(u)
	}

}

func AddTagRelation(u api.User) {
	for i := 0; i < 1; i++ {
		tag := strings.ToLower(u.Tags[i])
		q := `MATCH (u:User) WHERE u.username = {username} MATCH (n:TAG) WHERE n.value = "` + tag + `" CREATE (u)-[g:TAGGED]->(n) return n`
		st := app.prepareFakeStatement(q)
		api.ExecuteStatement(st, api.MapOf(u))
	}
}

func (app *App) insertFakeUser(u api.User) {
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
ilike: {ilike}, relation: {relation}, tags: {tags}})`
	//fmt.Println("Query == ", q)
	st := app.prepareFakeStatement(q)
	api.ExecuteStatement(st, api.MapOf(u))
	return
}

func (app *App) prepareFakeStatement(query string) bolt.Stmt {
	st, err := app.Neo.PrepareNeo(query)
	api.HandleError(err)
	return st
}
