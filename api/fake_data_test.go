package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestAddFakeData(t *testing.T) {
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
