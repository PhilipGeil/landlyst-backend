package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
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

func SendEmail() {
	from := mail.Address{Name: "Landlyst", Address: "phil2643@zbc.dk"}
	to := mail.Address{Name: "Philip", Address: "pgj@individualisterne.dk"}
	subject := "Bekræft oprettelse"
	host := "smtp.office365.com"

	// headers := make(map[string]string)
	// headers["From"] = from.String()
	// headers["To"] = to.String()
	// headers["Subject"] = subject
	// headers["MIME-Version"] = mime
	// headers["Content-Type"] = contentType
	// headers["Content-Transfer-Encoding"] = contentTransfer

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

	addressConfirmation, err := template.New("addressConfirmation").Parse(addressConfirmationTemplate)
	if err != nil {
		panic(err)
	}

	var body bytes.Buffer

	// t, err := template.New("template.html").ParseFiles("template.html")
	// if err != nil {
	// 	fmt.Println("Read file err")
	// }

	addressConfirmation.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Philip Jensen",
		Message: "This is a test message in a HTML template",
	})

	msg.WriteString(base64.StdEncoding.EncodeToString(body.Bytes()))

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

	if err = c.Auth(LoginAuth("phil2643@zbc.dk", "Alba2018")); err != nil {
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
