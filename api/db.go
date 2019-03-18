package api

import (
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

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

func (app *App) validToken(randomToken string) error {
	var u User
	err := app.Db.Get(&u, `SELECT id, access_lvl FROM users WHERE random_token=$1`, randomToken)
	if err != nil {
		return errors.New("Invalid Link")
	} else if u.AccessLvl == 1 {
		return errors.New("Email already validated")
	}
	_, err = app.Db.NamedExec(`UPDATE "public"."users" SET "access_lvl" = 1 WHERE "id" = :id`, u)
	return err
}

func tableOf(values string) string {
	return strings.Replace(values, ":", "", 99999)
}

func (app *App) insertUser(u User) {
	const vUsers = `(:username, :email, :lastname, :firstname, :password, :random_token, :img1, :img2, :img3, :img4, :img5, :biography, :birthday, :genre, :interest, :city, :zip, :country, :latitude, :longitude, :geo_allowed, :online, :rating, :access_lvl)`
	var query = `INSERT INTO public.users ` + tableOf(vUsers) + ` VALUES ` + vUsers
	_, err := app.Db.NamedExec(query, u)
	if err != nil {
		log.Fatalln(err)
	}
}
func (app *App) accountExist(rf registerForm) bool {
	u := User{}
	err := app.Db.Get(&u, `SELECT * FROM users WHERE username=$1 OR email=$2`, rf.Username, rf.Email)
	if u.Id != 0 || err == nil {
		return true
	} else {
		return false
	}
}

func (app *App) getUser(Username string) (u User, err error) {
	err = app.Db.Get(&u, `SELECT * FROM users WHERE username=$1`, Username)
	return
}

func (app *App) getBasicUser(Id int) (u User, err error) {
	const vBasicUser = `:username, :email, :lastname, :firstname, :img1, :img2, :img3, :img4, :img5, :biography, :birthday, :genre, :interest, :city, :zip, :country, :geo_allowed, :rating`
	err = app.Db.Get(&u, `SELECT `+tableOf(vBasicUser)+` FROM users WHERE id = $1`, Id)
	return
}

func (app *App) getBasicDates(Id int) ([]User, error) {
	const vDates = `:username, :email, :lastname, :firstname, :img1, :img2, :img3, :img4, :img5, :biography, :birthday, :genre, :interest, :city, :zip, :country, :geo_allowed, :rating`
	dates := []User{}
	err := app.Db.Select(&dates, `SELECT `+tableOf(vDates)+` FROM users ORDER BY id ASC`)
	return dates, err
}
