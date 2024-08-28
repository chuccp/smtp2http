package smtp

import (
	"errors"
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/util"
	"github.com/wneessen/go-mail"
	"os"
	"strings"
)

type Mail struct {
	Name string
	Mail string
}
type SMTP struct {
	Username string
	Mail     string
	Password string
	Host     string
	Port     int
}

type File struct {
	File     *os.File `json:"-"`
	Name     string   `json:"name"`
	FilePath string   `json:"filePath"`
}

type SendMsg struct {
	Subject       string
	BodyString    string
	ReceiveEmails []*Mail
	SendMail      *SMTP
	ff            []*File
}

func sendTestMsg(Smtp *SMTP) error {
	var sendMsg SendMsg
	sendMsg.SendMail = Smtp
	sendMsg.ReceiveEmails = []*Mail{{Name: Smtp.Username, Mail: Smtp.Mail}}
	sendMsg.Subject = "mail test"
	sendMsg.BodyString = " this is a test"
	return SendMail(&sendMsg)
}

func SendTestMsg(st *db.SMTP) error {
	if len(st.Username) == 0 {
		return errors.New("username cannot be empty")
	}
	if len(st.Password) == 0 {
		return errors.New("password cannot be empty")
	}
	if len(st.Host) == 0 {
		return errors.New("host cannot be empty")
	}
	return sendTestMsg(&SMTP{Username: st.Username, Mail: st.Mail, Password: st.Password, Host: st.Host, Port: st.Port})
}
func SendContentMsg(smtp *db.SMTP, mails []*db.Mail, subject, bodyString string) error {
	SMTP := &SMTP{Username: smtp.Username, Mail: smtp.Mail, Password: smtp.Password, Host: smtp.Host, Port: smtp.Port}
	receiveEmails := make([]*Mail, 0)
	for _, d := range mails {
		receiveEmails = append(receiveEmails, &Mail{Name: d.Name, Mail: d.Mail})
	}
	var sendMsg SendMsg
	sendMsg.SendMail = SMTP
	sendMsg.ReceiveEmails = receiveEmails
	sendMsg.Subject = subject
	sendMsg.BodyString = bodyString
	return SendMail(&sendMsg)
}
func SendFilesMsg(smtp *db.SMTP, mails []*db.Mail, files []*File, subject, bodyString string) error {
	SMTP := &SMTP{Username: smtp.Username, Mail: smtp.Mail, Password: smtp.Password, Host: smtp.Host, Port: smtp.Port}
	receiveEmails := make([]*Mail, 0)
	for _, d := range mails {
		receiveEmails = append(receiveEmails, &Mail{Name: d.Name, Mail: d.Mail})
	}
	var sendMsg SendMsg
	sendMsg.SendMail = SMTP
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
		if strings.Contains(err.Error(), "not found") && len(invalidMails) == 0 {
			mails := util.ExtractEmails(err.Error())
			err := SendMail(sendMsg, mails...)
			if err != nil {
				return err
			} else {
				return ToUserNotFoundError(mails)
			}
		}
		return err
	}
	return nil
}
