package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math"
	"testing"
)

func TestParseClaims(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTI3Mzk1MDYsImlkIjoxLCJsb2dnZWRJbiI6dHJ1ZSwidXNlcm5hbWUiOiJCTksifQ.SG7Hc7ukC46L0xLYq5FPJzrzGF4fbFE_h6qErIOsb9s"

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
