package stmp

import (
	"errors"
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/util"
	"github.com/wneessen/go-mail"
	"os"
	"strings"
)

type Mail struct {
	Name string
	Mail string
}
type STMP struct {
	Username string
	Mail     string
	Password string
	Host     string
	Port     int
}

type File struct {
	File *os.File
	Name string
}

type SendMsg struct {
	Subject       string
	BodyString    string
	ReceiveEmails []*Mail
	SendMail      *STMP
	ff            []*File
}

func sendTestMsg(STMP *STMP) error {
	var sendMsg SendMsg
	sendMsg.SendMail = STMP
	sendMsg.ReceiveEmails = []*Mail{{Name: STMP.Username, Mail: STMP.Mail}}
	sendMsg.Subject = "mail test"
	sendMsg.BodyString = " this is a test"
	return SendMail(&sendMsg)
}

func SendTestMsg(st *db.STMP) error {
	if len(st.Username) == 0 {
		return errors.New("username cannot be empty")
	}
	if len(st.Password) == 0 {
		return errors.New("password cannot be empty")
	}
	if len(st.Host) == 0 {
		return errors.New("host cannot be empty")
	}
	return sendTestMsg(&STMP{Username: st.Username, Mail: st.Mail, Password: st.Password, Host: st.Host, Port: st.Port})
}
func SendContentMsg(stmp *db.STMP, mails []*db.Mail, subject, bodyString string) error {
	STMP := &STMP{Username: stmp.Username, Mail: stmp.Mail, Password: stmp.Password, Host: stmp.Host, Port: stmp.Port}
	receiveEmails := make([]*Mail, 0)
	for _, d := range mails {
		receiveEmails = append(receiveEmails, &Mail{Name: d.Name, Mail: d.Mail})
	}
	var sendMsg SendMsg
	sendMsg.SendMail = STMP
	sendMsg.ReceiveEmails = receiveEmails
	sendMsg.Subject = subject
	sendMsg.BodyString = bodyString
	return SendMail(&sendMsg)
}
func SendFilesMsg(stmp *db.STMP, mails []*db.Mail, files []*File, subject, bodyString string) error {
	STMP := &STMP{Username: stmp.Username, Mail: stmp.Mail, Password: stmp.Password, Host: stmp.Host, Port: stmp.Port}
	receiveEmails := make([]*Mail, 0)
	for _, d := range mails {
		receiveEmails = append(receiveEmails, &Mail{Name: d.Name, Mail: d.Mail})
	}
	var sendMsg SendMsg
	sendMsg.SendMail = STMP
	sendMsg.ReceiveEmails = receiveEmails
	sendMsg.Subject = subject
	sendMsg.BodyString = bodyString
	for _, file := range files {
		sendMsg.ff = append(sendMsg.ff, file)
	}
	return SendMail(&sendMsg)
}

func SendMail(sendMsg *SendMsg, invalidMails ...string) error {

	msg := mail.NewMsg()
	if len(sendMsg.SendMail.Username) > 0 {
		err := msg.FromFormat(sendMsg.SendMail.Username, sendMsg.SendMail.Mail)
		if err != nil {
			return err
		}
	} else {
		if err := msg.From(sendMsg.SendMail.Mail); err != nil {
			return err
		}
	}

	for _, email := range sendMsg.ReceiveEmails {
		if !util.EqualsAnyIgnoreCase(email.Mail, invalidMails...) {
			if len(email.Name) > 0 {
				err := msg.AddToFormat(email.Name, email.Mail)
				if err != nil {
					continue
				}
			} else {
				msg.AddTo(email.Mail)
			}
		}
	}
	msg.Subject(sendMsg.Subject)
	msg.SetBodyString(mail.TypeTextPlain, sendMsg.BodyString)
	for _, f := range sendMsg.ff {
		err := msg.AttachReader(f.Name, f.File, mail.WithFileName(f.Name))
		if err != nil {
			continue
		}
	}
	c, err := mail.NewClient(sendMsg.SendMail.Host, mail.WithPort(sendMsg.SendMail.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(sendMsg.SendMail.Username), mail.WithPassword(sendMsg.SendMail.Password))
	if err != nil {
		return err
	}
	if err := c.DialAndSend(msg); err != nil {
		if strings.Contains(err.Error(), "User not found") {
			mails := util.ExtractEmails(err.Error())
			return SendMail(sendMsg, mails...)
		}
		return err
	}
	return nil
}
