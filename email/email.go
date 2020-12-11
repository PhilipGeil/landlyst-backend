package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

	"github.com/PhilipGeil/landlyst-backend/api/resources"
)

type loginAuth struct {
	username, password string
}

type Mail interface {
	WriteMail(w io.Writer) error
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func SendVerifyEmail(uuid, email, fname string) {
	subject := "Bekræft oprettelse"
	var body bytes.Buffer

	tmpl := template.Must(template.ParseFiles("C:\\Users\\phil2643\\development\\landlyst\\api-server\\email\\template.html"))

	if err := tmpl.Execute(&body, struct {
		Name string
		Link string
	}{
		Name: "Philip Jensen",
		Link: "http://localhost:8080/api/verify/" + uuid,
	}); err != nil {
		log.Fatalln(err)
	}

	SendEmail(base64.StdEncoding.EncodeToString(body.Bytes()), fname, email, subject)
}

func SendConfirmEmail(res resources.ReservationResponse) {
	subject := "Ordre bekræftelse"

	var body bytes.Buffer

	tmpl := template.Must(template.ParseFiles("C:\\Users\\phil2643\\development\\landlyst\\api-server\\email\\confirmed.html"))

	if err := tmpl.Execute(&body, struct {
		User      resources.Customer
		Link      string
		StartDate string
		EndDate   string
	}{
		User:      res.Reservation.Customer,
		Link:      "http://localhost:3000/login",
		StartDate: res.Reservation.Dates.StartDate.Format("2006-01-02 15:04:05"),
		EndDate:   res.Reservation.Dates.EndDate.Format("2006-01-02 15:04:05"),
	}); err != nil {
		log.Fatalln(err)
	}

	SendEmail(base64.StdEncoding.EncodeToString(body.Bytes()), res.Reservation.Customer.FName, res.Reservation.Customer.Email, subject)
}

func SendEmail(s, fname, email, subject string) {
	from := mail.Address{Name: "Landlyst Kro og Hotel", Address: "phil2643@zbc.dk"}
	to := mail.Address{Name: fname, Address: email}
	host := "smtp.office365.com"
	authEmail := os.Getenv("EMAIL_AUTH_EMAIL")
	authPass := os.Getenv("EMAIL_AUTH_PASS")

	var msg strings.Builder
	msg.WriteString("From: ")
	msg.WriteString(from.String())
	msg.WriteString("\r\n")
	msg.WriteString("To: ")
	msg.WriteString(to.String())
	msg.WriteString("\r\n")
	msg.WriteString("Subject: ")
	msg.WriteString(mime.QEncoding.Encode("utf-8", subject))
	msg.WriteString("\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: base64\r\n")
	msg.WriteString("\r\n")

	msg.WriteString(s)

	tlsconfig := &tls.Config{
		ServerName: host,
	}

	conn, err := net.Dial("tcp", "smtp.office365.com:587")
	if err != nil {
		fmt.Println("tls.Dial Error: ", err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Println("smtp.NewClient Error: ", err)
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		fmt.Println("c.StartTLS Error: ", err)
		return
	}

	if err = c.Auth(LoginAuth(authEmail, authPass)); err != nil {
		fmt.Println("c.Auth Error: ", err)
		return
	}

	if err = c.Mail(from.Address); err != nil {
		fmt.Println("c.Mail Error: ", err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		fmt.Println("c.Rcpt Error: ", err)
	}

	w, err := c.Data()
	if err != nil {
		fmt.Println("c.Data Error: ", err)
	}

	_, err = w.Write([]byte(msg.String()))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = w.Close()
	if err != nil {
		fmt.Println("reader Error: ", err)
	}

	c.Quit()
}
