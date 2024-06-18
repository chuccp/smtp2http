package stmp

import (
	"errors"
	"github.com/chuccp/d-mail/db"
	"github.com/wneessen/go-mail"
	"os"
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

type SendMsg struct {
	Subject       string
	BodyString    string
	ReceiveEmails []*Mail
	SendMail      *STMP
	ff            []*os.File
}

func (sendMsg *SendMsg) GetFromMail() string {
	return sendMsg.SendMail.Mail
}
func (sendMsg *SendMsg) GetToMail() []string {
	datas := make([]string, 0)
	for _, email := range sendMsg.ReceiveEmails {
		datas = append(datas, email.Mail)
	}
	return datas
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

func SendMail(sendMsg *SendMsg) error {

	msg := mail.NewMsg()

	if err := msg.From(sendMsg.GetFromMail()); err != nil {
		return err
	}
	if err := msg.To(sendMsg.GetToMail()...); err != nil {
		return err
	}

	msg.Subject(sendMsg.Subject)
	msg.SetBodyString(mail.TypeTextPlain, sendMsg.BodyString)

	var files []*mail.File
	for _, f := range sendMsg.ff {
		err := msg.AttachReader(f.Name(), f, mail.WithFileName(f.Name()))
		if err != nil {
			continue
		}
	}
	msg.SetAttachements(files)
	c, err := mail.NewClient(sendMsg.SendMail.Host, mail.WithPort(sendMsg.SendMail.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(sendMsg.SendMail.Username), mail.WithPassword(sendMsg.SendMail.Password))
	if err != nil {
		return err
	}
	if err := c.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
