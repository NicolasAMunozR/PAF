package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strconv"

	"gopkg.in/mail.v2"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

var (
	email, pass, server, port string
)

func Setup() {
	email = os.Getenv("EMAIL_DIR")
	pass = os.Getenv("EMAIL_PASS")
	server = os.Getenv("EMAIL_SERVER")
	port = os.Getenv("EMAIL_PORT")
}

func NewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() error {
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%s", server, port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", email, pass, server), email, r.to, []byte(body)); err != nil {
		return err
	}
	return nil
}

func (r *Request) Send(templateName string, items interface{}) error {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		return err
	}
	if err := r.sendMail(); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *Request) SendMailSkipTLS(templateName string, items interface{}) error {

	err := r.parseTemplate(templateName, items)
	if err != nil {
		return err
	}

	m := mail.NewMessage()

	m.SetHeader("From", email)
	m.SetHeader("To", r.to[0]) //se pueden colocar mas mail separados por coma : m.SetHeader("To", "email1@email.com","email2@email.com")
	//m.SetAddressHeader("Cc", "manuel.manriquez.lopez@gmail.com", "Manuel")
	m.SetHeader("Subject", r.subject)
	m.SetBody("text/html", r.body)
	//m.Attach("/home/Alex/lolcat.jpg")

	//d := mail.NewDialer("smtp.example.com", 587, "user", "123456")
	//d.StartTLSPolicy = mail.MandatoryStartTLS
	var puerto int
	puerto, err = strconv.Atoi(port)
	if err != nil {
		fmt.Printf("No se pudo mandar el email: %s", err.Error())
		return err
	}
	d := mail.NewDialer(server, puerto, email, pass)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("No se pudo mandar el email2: %s", err.Error())
		return err
	}

	return nil
}
