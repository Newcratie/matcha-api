package api

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

func validateUser(rf registerForm) (User, error) {

	// returns nil or ValidationErrors ( []FieldError )
	err := app.checkRegister(rf)
	if err != nil {
		return User{}, err
	}
	err = app.validate.Struct(rf)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return User{}, err
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			fmt.Println(err.StructField())     // by passing alt name to ReportError like below
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		errs := errors.New("")
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				if err.Tag() == "alphanumunicode" {
					errs = errors.New(errs.Error() + "Username must contain only alpha numeric characters!")
				} else {
					errs = errors.New(errs.Error() + "Username should be between 6 and 20 characters long\n")
				}
				break
			case "Email":
				errs = errors.New(errs.Error() + "Invalid email\n")
				break
			case "Password":
				errs = errors.New(errs.Error() + "Invalid password\n")
				break
			case "Confirm":
				errs = errors.New(errs.Error() + "Passwords don't match\n")
				break
			case "Firstname":
				errs = errors.New(errs.Error() + "Invalid Firstname\n")
				break
			case "Birthday":
				errs = errors.New(errs.Error() + "Invalid birthaday\n")
				break
			}
		}

		// from here you can create your own error messages in whatever language you wish
		return User{}, errs
	}

	// save user to database
	user := NewUser(rf)
	return user, nil
}

func validateVariable() {
	myEmail := "joeybloggs.gmail.com"

	errs := app.validate.Var(myEmail, "required,email")

	if errs != nil {
		fmt.Println(errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
		return
	}

	// email ok, move on
}
