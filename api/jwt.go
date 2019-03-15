package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var key = []byte(hashKey)

func (user User) GenerateJwt() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"loggedIn": true,
		"id":       user.Id,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(key)
	if err != nil {
		fmt.Println("something get wrong with jwt: ", err.Error())
		return "", err
	}
	return tokenString, nil
}

func checkJwt(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(hashKey), nil
	})

	if token.Valid {
		return true
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
			return false
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
			return false
		} else {
			fmt.Println("Couldn't handle this token:", err)
			return false
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
		return false
	}
	return true
}
