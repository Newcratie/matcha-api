package api

import (
	"github.com/Newcratie/matcha-api/api/logprint"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

func TestFetchUser(t *testing.T) {
	app.newApp()
	go app.routerAPI()
	app.Db = dbConnect()
	app.fetchUsers()
	logprint.Title("Users")
	logprint.PrettyPrint(app.Users)
	app.validate = validator.New()
	if len(app.Users) == 0 {
		t.Error("Api Cant fetch DB")
	}
	logprint.End()
}
