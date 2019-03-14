package api

import (
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

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

func (app *App) insertUser(u User) {
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
