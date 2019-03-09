package api

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

type validationField struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type msgs struct {
	UserAlpha       validationField `json:"user_alpha"`
	UserLen         validationField `json:"user_len"`
	UserExist       validationField `json:"user_exist"`
	EmailInvalid    validationField `json:"email_invalid"`
	PasswordInvalid validationField `json:"password_invalid"`
	PasswordMatch   validationField `json:"password_match"`
	Firstname       validationField `json:"firstname"`
	Lastname        validationField `json:"lastname"`
	Birthday        validationField `json:"birthday"`
}
type validationRes struct {
	Valid bool `json:"valid"`
	Fail  bool `json:"fail"`
	Errs  msgs `json:"errs"`
}

func (res *validationRes) failure() {
	res.Valid = false
	res.Fail = true
}

func validateUser(rf registerForm) (User, validationRes) {
	es := msgs{
		validationField{true, ""},
		validationField{true, ""},
		validationField{true, ""},
		validationField{true, ""},
		validationField{true, ""},
		validationField{true, ""},
		validationField{true, ""},
		validationField{true, ""},
		validationField{true, ""},
	}
	res := validationRes{
		true,
		false,
		es,
	}
	exist := app.accountExist(rf)
	if exist {
		res.Errs.UserExist.Status = false
		res.Errs.UserExist.Message = "Username or Email already exist"

	}
	err := app.validate.Struct(rf)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			res.failure()
			return User{}, res
		}

		for _, err := range err.(validator.ValidationErrors) {
			//fmt.Println(err.Namespace())
			//fmt.Println(err.Field())
			//fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			//fmt.Println(err.StructField())     // by passing alt name to ReportError like below
			fmt.Println(err.Tag())
			//fmt.Println(err.ActualTag())
			//fmt.Println(err.Kind())
			//fmt.Println(err.Type())
			//fmt.Println(err.Value())
			//fmt.Println(err.Param())
			//fmt.Println()
		}

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				if err.Tag() == "alphanumunicode" {
					res.Errs.UserAlpha.Status = false
					res.Errs.UserAlpha.Message = "Username must contain only alpha numeric characters!"
				} else {
					res.Errs.UserAlpha.Status = false
					res.Errs.UserAlpha.Message = "Username must contain between 6 and 20 characters"
				}
				break
			case "Email":
				res.Errs.EmailInvalid.Status = false
				res.Errs.EmailInvalid.Message = "Invalid Email"
				break
			case "Password":
				res.Errs.PasswordInvalid.Status = false
				res.Errs.PasswordInvalid.Message = "Invalid password"
				break
			case "Confirm":
				res.Errs.PasswordMatch.Status = false
				res.Errs.PasswordMatch.Message = "Password don't match"
				break
			case "Firstname":
				res.Errs.Firstname.Status = false
				res.Errs.Firstname.Message = "Invalid Firstname"
				break
			case "Lastname":
				res.Errs.Firstname.Status = false
				res.Errs.Firstname.Message = "Invalid Lastname"
				break
			case "Birthday":
				res.Errs.Birthday.Status = false
				res.Errs.Birthday.Message = "Invalid birthday"
				break
			}
		}

		// from here you can create your own error messages in whatever language you wish
		res.failure()
		return User{}, res
	}
	var u User

	u.Username = rf.Username
	u.Email = rf.Email
	u.LastName = rf.Lastname
	u.FirstName = rf.Firstname
	u.Password = rf.Username
	u.Birthday, _ = parseTime(rf.Birthday)
	return u, res
}

func (app *App) accountExist(rf registerForm) bool {
	u := User{}
	err := app.Db.Get(&u, `SELECT * FROM users WHERE username=$1 OR email=$2`, rf.Username, rf.Email)
	fmt.Println("checkRegister: >>", u.Id, "<<", err)
	if u.Id != 0 {
		return true
	} else {
		return false
	}
}
