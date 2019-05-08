package api

import (
	"github.com/Newcratie/matcha-api/api/hash"
	"github.com/brianvoe/gofakeit"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func newRandomMale() User {
	var f *gofakeit.PersonInfo
	max := 1
	f = gofakeit.Person()
	interest := make([]string, 3)
	interest[0] = "bi"
	interest[1] = "hetero"
	interest[2] = "homo"
	tagtab := make([]string, 4)
	//tagg := []Tag
	for i := 0; i < max; i++ {
		tagtab[i] = gofakeit.Color()
		//tagg.
	}

	Latitude, _ := gofakeit.LatitudeInRange(42.490627, 50.264989)
	Longitude, _ := gofakeit.LongitudeInRange(-3.396493, 9.517944)
	return User{Username: gofakeit.Username(),
		Password:  hash.Encrypt(hashKey, "'"),
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

func newRandomFemale() User {
	Latitude, _ := gofakeit.LatitudeInRange(42.490627, 50.264989)
	Longitude, _ := gofakeit.LongitudeInRange(-3.396493, 9.517944)
	var f *gofakeit.PersonInfo
	max := 1
	f = gofakeit.Person()
	interest := make([]string, 3)
	interest[0] = "bi"
	interest[1] = "hetero"
	interest[2] = "homo"
	tagtab := make([]string, 4)
	for i := 0; i < max; i++ {
		tagtab[i] = gofakeit.Color()
	}
	return User{Username: gofakeit.Username(),
		Password:  hash.Encrypt(hashKey, "'"),
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

func TestAddFakeData(t *testing.T) {
	const max = 50
	host := os.Getenv("NEO_HOST")
	app.Db, _ = bolt.NewDriverPool("bolt://neo4j:secret@"+host+":7687", 1000)
	app.Neo, _ = app.Db.OpenPool()
	for i := 0; i < max; i++ {
		s := gofakeit.Color()
		s = strings.ToLower(s)
		app.Neo.QueryNeoAll(`MERGE (t:TAG {key: "`+s+`", text: "#`+strings.Title(s)+`", value: "`+s+`"}) `, nil)
	}
	for i := 0; i < max; i++ {
		u := newRandomMale()
		app.insertUser(u)
		AddTagRelation(u)
		u = newRandomFemale()
		app.insertUser(u)
		AddTagRelation(u)
	}

}

func AddTagRelation(u User) {
	for i := 0; i < 4; i++ {
		tag := strings.ToLower(u.Tags[i])
		q := `MATCH (u:User) WHERE u.username = {username} MATCH (n:TAG) WHERE n.value = "` + tag + `" CREATE (u)-[g:TAGGED]->(n) return n`
		st := app.prepareStatement(q)
		executeStatement(st, MapOf(u))
	}
}
