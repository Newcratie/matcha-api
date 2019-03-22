package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"math"
	"os"
	"testing"
)

func testNeo(t *testing.T) {

	driver := bolt.NewDriver()
	host := os.Getenv("NEO_HOST")
	app.Neo, _ = driver.OpenNeo("bolt://neo4j:secret@" + host + ":7687")
	data, _, _, _ := app.Neo.QueryNeoAll(`MATCH (n:User{random_token : "eba2beb1bc9c69315bb36946c7adfe40f0ec9706c33d63e23201c6a6bc100345"}) SET n.access_lvl = 1 RETURN n`, nil)
	if len(data) == 0 {
	} else if data[0][0].(graph.Node).Properties["access_lvl"] == int64(1) {
	} else {
		fmt.Printf("type %T\n", data[0][0].(graph.Node).Properties["access_lvl"])
	}
}

func TestParseClaims(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTMyMDc3NzYsImlkIjowLCJsb2dnZWRJbiI6dHJ1ZSwidXNlcm5hbWUiOiJCTksifQ.8zXBRvkysk0SfZxjU8-ThIwo_DXUHswsr5uferPAmyc"

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(hashKey), nil
	})
	if err != nil {

	} else if checkJwt(tokenString) {
		id := int(math.Round(claims["id"].(float64)))
		fmt.Printf(" Type: %d\n", id)
	}
	//for key, val := range claims {
	//	fmt.Printf("Key: %v, value: %v\n", key, val)
	//}
}
