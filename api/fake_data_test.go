package api

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func newRandomUser() User {
	var f *gofakeit.PersonInfo
	f = gofakeit.Person()
	interest := make([]string, 3)
	interest[0] = "bi"
	interest[1] = "hetero"
	interest[2] = "homo"
	return User{Username: gofakeit.Username(),
		Password:  "fakepass",
		FirstName: f.FirstName,
		LastName:  f.LastName,
		Email:     gofakeit.Email(),
		Img1:      "https://randomuser.me/api/portraits/men/" + string(gofakeit.Number(1, 45)) + ".jpg",
		Img2:      "https://randomuser.me/api/portraits/men/" + string(gofakeit.Number(1, 45)) + ".jpg",
		Img3:      gofakeit.ImageURL(300, 300),
		Img4:      gofakeit.ImageURL(300, 300),
		Img5:      gofakeit.ImageURL(300, 300),
		Biography: gofakeit.Paragraph(2, 30, 12, " "),
		Birthday: gofakeit.DateRange(time.Date(1900, 01, 01, 00, 00, 00, 00, time.UTC),
			time.Date(2000, 01, 01, 00, 00, 00, 00, time.UTC)),
		Genre:      f.Gender,
		Interest:   gofakeit.RandString(interest),
		AccessLvl:  1,
		Online:     gofakeit.Bool(),
		Rating:     gofakeit.Float32Range(0, 10),
		City:       gofakeit.City(),
		Zip:        gofakeit.Zip(),
		Country:    gofakeit.Country(),
		Latitude:   gofakeit.Latitude(),
		Longitude:  gofakeit.Longitude(),
		GeoAllowed: gofakeit.Bool(),
		CreatedAt: gofakeit.DateRange(time.Date(1900, 01, 01, 00, 00, 00, 00, time.Local),
			time.Date(2017, 01, 01, 00, 00, 00, 00, time.Local)),
	}
}

func TestAddFakeData(t *testing.T) {
	driver := bolt.NewDriver()
	host := os.Getenv("NEO_HOST")
	app.Neo, _ = driver.OpenNeo("bolt://neo4j:secret@" + host + ":7687")
	for i := 0; i < 100; i++ {
		u := newRandomUser()
		app.insertUser(u)
	}
}

func estAddFakeData(t *testing.T) {
	app.newApp()
	app.Db = dbConnect()
	data, err := ioutil.ReadFile("./tests/fake_data.json")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	var u []User
	err = json.Unmarshal(data, &u)
	if err != nil {
		panic(err)
	}
	for _, user := range u {
		app.insertUser(user)
	}
}
