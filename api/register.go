package api

import (
	"crypto/rand"
	"fmt"
	"github.com/Newcratie/matcha-api/api/hash"
	"regexp"
)

const min = 2
const max = 20

func validateUser(rf registerForm) (User, validationResponse) {
	res := validationResponse{true, false, make([]ErrorField, 0)}
	if app.accountExist(rf) {
		res.failure("Account already exist")
	}
	if rf.Password != rf.Confirm {
		res.failure("Passwords don't match")
	}
	if !emailIsValid(rf.Email) {
		res.failure("Invalid email")
	}
	if len(rf.Username) < min || len(rf.Username) > max {
		res.failure("Username must contain between " + string(min) + " and " + string(max) + " characters")
	}
	if len(rf.Lastname) < min || len(rf.Lastname) > max {
		res.failure("Lastname must contain between " + string(min) + " and " + string(max) + " characters")
	}
	if len(rf.Firstname) < min || len(rf.Firstname) > max {
		res.failure("Firstname must contain between " + string(min) + " and " + string(max) + " characters")
	}

	if res.Valid {
		return user(rf), res
	} else {
		return User{}, res
	}
}

func user(rf registerForm) (u User) {
	u.Username = rf.Username
	u.Email = rf.Email
	u.Password = hash.Encrypt(hashKey, rf.Password)
	u.LastName = rf.Lastname
	u.FirstName = rf.Firstname
	u.Birthday = rf.Birthday
	u.RandomToken = newToken()
	sendToken(u.Username, u.Email, u.RandomToken)
	return
}

func (res *validationResponse) failure(msg string) {
	res.Valid = false
	res.Fail = true
	res.Errs = append(res.Errs, ErrorField{false, msg})
}

func newToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func emailIsValid(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
}
