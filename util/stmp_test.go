package util

import (
	"github.com/wneessen/go-mail"
	"log"
	"net/smtp"
	"testing"
)

func TestAAA(t *testing.T) {

	from := "chuccp@163.com"
	hostname := "smtp.163.com"
	auth := smtp.PlainAuth("", "chuccp@163.com", "VBHGROQFVIEPMAMZ", hostname)
	var recipients []string = []string{"chuccp@163.com"}
	msg := []byte("dummy message")
	err := smtp.SendMail(hostname+":25", auth, from, recipients, msg)
	if err != nil {
		log.Fatal(err)
	}

}
func TestBBBB(t *testing.T) {

	m := mail.NewMsg()
	if err := m.From("chuccp@163.com"); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := m.To("chuccp@163.com"); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}

	m.Subject("This is my first mail with go-mail!")
	m.SetBodyString(mail.TypeTextPlain, "Do you like this mail? I certainly do!")
	c, err := mail.NewClient("smtp.163.com", mail.WithPort(25), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername("chuccp"), mail.WithPassword("VBHGROQFVIEPMAMZ"))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	if err := c.DialAndSend(m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}

}
