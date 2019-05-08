package api

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/Newcratie/matcha-api/api/logprint"
	"html/template"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
)

func SendEmail(Title, username, email, link string) error {
	logprint.Title("send Token")
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpLogin := "camagru4422@gmail.com"
	smtpPasswd := "42istheanswer"

	templateData := struct {
		Name string
		URL  string
	}{
		Name: username,
		URL:  "http://localhost:8080/valid_email?token=",
	}

	from := mail.Address{"", smtpLogin}
	to := mail.Address{username, email}
	title := Title
	body, err := ParseTemplate(link, templateData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	useTls := false
	useStartTls := true

	header := make(map[string]string)
	header["From"] = "matcha@42.fr"
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
		return err
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
		return err
	}

	hasStartTLS, _ := client.Extension("STARTTLS")
	if useStartTls && hasStartTLS {
		fmt.Println("STARTTLS ...")
		if err = client.StartTLS(tlsconfig); err != nil {
			fmt.Println(err)
			return err
		}
	}

	auth := smtp.PlainAuth(
		"",
		smtpLogin,
		smtpPasswd,
		smtpHost,
	)

	if ok, _ := client.Extension("AUTH"); ok {
		if err := client.Auth(auth); err != nil {
			fmt.Printf("Error during AUTH %s\n", err)
			return err
		}
	}
	fmt.Println("AUTH done")

	if err := client.Mail(from.Address); err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	fmt.Println("FROM done")

	if err := client.Rcpt(to.Address); err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Address: %s\n", to.Address)
		return err
	}
	fmt.Println("TO done")

	w, err := client.Data()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	err = w.Close()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	client.Quit()
	return nil
}

func SendEmailValidation(username, email, token string) error {
	logprint.Title("send Token")
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpLogin := "camagru4422@gmail.com"
	smtpPasswd := "42istheanswer"

	templateData := struct {
		Name string
		URL  string
	}{
		Name: username,
		URL:  "http://localhost:8080/valid_email?token=" + token,
	}

	from := mail.Address{"", smtpLogin}
	to := mail.Address{username, email}
	title := "Validate your address"
	body, err := ParseTemplate("./api/utils/confirm_email.html", templateData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	useTls := false
	useStartTls := true

	header := make(map[string]string)
	header["From"] = "matcha@42.fr"
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
		return err
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
		return err
	}

	hasStartTLS, _ := client.Extension("STARTTLS")
	if useStartTls && hasStartTLS {
		fmt.Println("STARTTLS ...")
		if err = client.StartTLS(tlsconfig); err != nil {
			fmt.Println(err)
			return err
		}
	}

	auth := smtp.PlainAuth(
		"",
		smtpLogin,
		smtpPasswd,
		smtpHost,
	)

	if ok, _ := client.Extension("AUTH"); ok {
		if err := client.Auth(auth); err != nil {
			fmt.Printf("Error during AUTH %s\n", err)
			return err
		}
	}
	fmt.Println("AUTH done")

	if err := client.Mail(from.Address); err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	fmt.Println("FROM done")

	if err := client.Rcpt(to.Address); err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Address: %s\n", to.Address)
		return err
	}
	fmt.Println("TO done")

	w, err := client.Data()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	err = w.Close()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	client.Quit()
	return nil
}

func SendEmailPasswordForgot(username, email, token string) error {
	logprint.Title("send Token")
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpLogin := "camagru4422@gmail.com"
	smtpPasswd := "42istheanswer"

	templateData := struct {
		Name string
		URL  string
	}{
		Name: username,
		URL:  "http://localhost:8080/reset-password?reset_token=" + token,
	}

	from := mail.Address{"", smtpLogin}
	to := mail.Address{username, email}
	title := "Password Changed"
	body, err := ParseTemplate("./api/utils/pass_change.html", templateData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	useTls := false
	useStartTls := true

	header := make(map[string]string)
	header["From"] = "matcha@42.fr"
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
		return err
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
		return err
	}

	hasStartTLS, _ := client.Extension("STARTTLS")
	if useStartTls && hasStartTLS {
		fmt.Println("STARTTLS ...")
		if err = client.StartTLS(tlsconfig); err != nil {
			fmt.Println(err)
			return err
		}
	}

	auth := smtp.PlainAuth(
		"",
		smtpLogin,
		smtpPasswd,
		smtpHost,
	)

	if ok, _ := client.Extension("AUTH"); ok {
		if err := client.Auth(auth); err != nil {
			fmt.Printf("Error during AUTH %s\n", err)
			return err
		}
	}
	fmt.Println("AUTH done")

	if err := client.Mail(from.Address); err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	fmt.Println("FROM done")

	if err := client.Rcpt(to.Address); err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Address: %s\n", to.Address)
		return err
	}
	fmt.Println("TO done")

	w, err := client.Data()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	err = w.Close()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	client.Quit()
	return nil
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
