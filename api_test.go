package api

import (
	"encoding/json"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func TestFetchUser(t *testing.T) {
	app.newApp()
	go app.routerAPI()
	app.Db = dbConnect()
	app.fetchUsers()
	app.validate = validator.New()
	PrettyPrint(app.Users)
	if len(app.Users) == 0 {
		t.Error("Api Cant fetch DB")
	}
}
