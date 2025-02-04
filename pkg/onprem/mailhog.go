package onprem

import (
	"os"
	"strconv"

	"github.com/shahzodshafizod/gocloud/pkg"
	"gopkg.in/gomail.v2"

	"context"

	"github.com/pkg/errors"
)

type email struct {
	username   string
	password   string
	from       string
	smtpserver string
	port       int
}

func NewEmail() (pkg.Email, error) {
	port, err := strconv.Atoi(os.Getenv("MAILHOG_PORT"))
	if err != nil {
		return nil, errors.Wrap(err, "strconv.Atoi")
	}

	client := &email{
		username:   os.Getenv("MAILHOG_USERNAME"),
		password:   os.Getenv("MAILHOG_PASSWORD"),
		from:       os.Getenv("MAILHOG_FROM"),
		smtpserver: os.Getenv("MAILHOG_SMTP_SERVER"),
		port:       port,
	}

	return client, nil
}

func (e *email) Send(ctx context.Context, emailTo string, subject string, body string) error {

	msg := gomail.NewMessage()
	msg.SetHeader("From", e.from)
	msg.SetHeader("To", emailTo)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dialer := gomail.NewDialer(e.smtpserver, e.port, e.username, e.password)
	err := dialer.DialAndSend(msg)
	if err != nil {
		return errors.Wrap(err, "dialer.DialAndSend")
	}

	return nil
}

// // https://gist.github.com/guinso/8405a991d8a095b01427b9ea83934d67

// func (e *email) Send2(ctx context.Context, to string, subject string, body string) error {

// 	tlsConfig := tls.Config{
// 		ServerName:         e.smtpserver,
// 		InsecureSkipVerify: true,
// 	}

// 	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", e.smtpserver, e.port), &tlsConfig)
// 	if err != nil {
// 		return errors.Wrap(err, "tls.Dial")
// 	}
// 	// defer conn.Close()

// 	client, err := smtp.NewClient(conn, e.smtpserver)
// 	if err != nil {
// 		return errors.Wrap(err, "smtp.NewClient")
// 	}
// 	// defer client.Close()

// 	auth := smtp.PlainAuth("", e.username, e.password, e.smtpserver)

// 	err = client.Auth(auth)
// 	if err != nil {
// 		return errors.Wrap(err, "client.Auth")
// 	}

// 	err = client.Mail(e.username)
// 	if err != nil {
// 		return errors.Wrap(err, "client.Mail")
// 	}

// 	// for _, to := range tos {
// 	err = client.Rcpt(to)
// 	if err != nil {
// 		return errors.Wrap(err, "client.Rcpt")
// 	}
// 	// }

// 	var msg string
// 	msg += fmt.Sprintf("From: %s\r\n", e.username)
// 	msg += fmt.Sprintf("To: %s\r\n", to) // strings.Join(tos, ","))
// 	msg += fmt.Sprintf("Subject: %s\r\n", subject)
// 	msg += "Content-Type: text/html; charset=utf-8\r\n"
// 	msg += "\r\n"
// 	msg += body

// 	writer, err := client.Data()
// 	if err != nil {
// 		return errors.Wrap(err, "client.Data")
// 	}

// 	_, err = writer.Write([]byte(msg))
// 	if err != nil {
// 		return errors.Wrap(err, "writer.Write")
// 	}

// 	err = writer.Close()
// 	if err != nil {
// 		return errors.Wrap(err, "writer.Close")
// 	}

// 	return nil
// }
