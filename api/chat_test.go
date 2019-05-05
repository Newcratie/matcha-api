package api

import (
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"os"
	"testing"
)

func TestWebSocket(t *testing.T) {
	driver := bolt.NewDriver()
	host := os.Getenv("NEO_HOST")
	app.Neo, _ = driver.OpenNeo("bolt://neo4j:secret@" + host + ":7687")
}

func TestCreateLike(t *testing.T) {

	//claims := jwt.MapClaims{}
	//valid, err := ValidateToken(c, &claims)

	var M Match
	M.IdFrom = 201
	prin("ID_TO ==>", M.IdFrom, "|")
	M.IdTo = 6
	prin("ID_TO ==>", M.IdTo, "|")
	M.Action = "LIKE"
	prin("ACTION ==>> ", M.Action, "|")
	if valid, err := app.dbMatchs(M); valid != true || err != nil {
		prin("VALID ==>", valid, "ERROR ==>", err, "|")
	}
	//fmt.Println("jwt error: ", err)
	//c.JSON(201, gin.H{"err": err.Error()})

}
