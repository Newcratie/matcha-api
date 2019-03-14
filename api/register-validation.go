package api

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/Newcratie/matcha-api/api/hash"
	"html/template"
	"net"
	"net/mail"
	"net/smtp"
	"regexp"
	"strconv"
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
	return
}

func (res *validationResponse) failure(msg string) {
	res.Valid = false
	res.Fail = true
	res.Errs = append(res.Errs, ErrorField{false, msg})
}

func (app *App) accountExist(rf registerForm) bool {
	u := User{}
	err := app.Db.Get(&u, `SELECT * FROM users WHERE username=$1 OR email=$2`, rf.Username, rf.Email)
	if u.Id != 0 || err == nil {
		return true
	} else {
		return false
	}
}

func sendToken(username, email, token string) {
	smtpHost := "mail.gmx.com"
	smtpPort := 587
	smtpLogin := "matcha42@gmx.com"
	smtpPasswd := "42born2code"

	templateData := struct {
		Name string
		URL  string
	}{
		Name: username,
		URL:  "http://localhost/auth/token/" + token,
	}

	from := mail.Address{"", smtpLogin}
	to := mail.Address{username, email}
	title := "Validate your address"
	body, err := ParseTemplate("./utils/confirm_email.html", templateData)
	if err != nil {
		fmt.Println(err)
	}

	useTls := false
	useStartTls := true

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = title
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	conn, err := net.Dial("tcp", smtpHost+":"+strconv.Itoa(smtpPort))
	if err != nil {
		fmt.Println(err)
		return
	}

	// TLS
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	if useTls {
		conn = tls.Client(conn, tlsconfig)
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Println(err)
		return
	}

	hasStartTLS, _ := client.Extension("STARTTLS")
	if useStartTls && hasStartTLS {
		fmt.Println("STARTTLS ...")
		if err = client.StartTLS(tlsconfig); err != nil {
			fmt.Println(err)
			return
		}
	}

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		smtpLogin,
		smtpPasswd,
		smtpHost,
	)

	if ok, _ := client.Extension("AUTH"); ok {
		if err := client.Auth(auth); err != nil {
			fmt.Printf("Error during AUTH %s\n", err)
			return
		}
	}
	fmt.Println("AUTH done")

	if err := client.Mail(from.Address); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Println("FROM done")

	if err := client.Rcpt(to.Address); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Println("TO done")

	w, err := client.Data()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	err = w.Close()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	client.Quit()
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, err
}

func newToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func emailIsValid(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
}
