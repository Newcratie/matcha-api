package api

import (
	"errors"
	"fmt"
	"github.com/fatih/structs"
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
	m := structs.Map(u)
	//var query = `CREATE (np:Person { Username: {Username}, Password: {Password}})`
	fmt.Println(m)
	result, _ := app.Neo.ExecNeo("CREATE (n:NODE {name: {Username}})", m)
	numResult, _ := result.RowsAffected()
	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	//_, err := app.Neo.ExecNeo(query, structs.Map(User{Username: "test2"}))
}
func (app *App) accountExist(rf registerForm) bool {
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:Person) WHERE n.username = "`+rf.Username+`" RETURN n`, nil)
	if len(data) == 0 {
		return false
	} else {
		return true
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
