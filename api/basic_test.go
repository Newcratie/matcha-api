package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"math"
	"os"
	"testing"
)

func TestNeo(t *testing.T) {

	driver := bolt.NewDriver()
	host := os.Getenv("NEO_HOST")
	app.Neo, _ = driver.OpenNeo("bolt://neo4j:secret@" + host + ":7687")
	u, err := app.getBasicDates(81)
	fmt.Println(err, u)
}

func testParseClaims(t *testing.T) {
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
