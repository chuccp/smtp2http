package service

import (
	"bytes"
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/util"
	"github.com/chuccp/smtp2http/web"
)

type Token struct {
	context *core.Context
}

func NewToken(context *core.Context) *Token {
	return &Token{context: context}
}
func (token *Token) GetPage(page *web.Page) (any, error) {
	p, err := token.context.GetDb().GetTokenModel().Page(page)
	if err != nil {
		return nil, err
	}
	token.supplement(p.List...)
	return p, nil
}
func (token *Token) GetOne(id int) (*db.Token, error) {
	one, err := token.context.GetDb().GetTokenModel().GetOne(uint(id))
	if err != nil {
		return nil, err
	}
	token.supplement(one)
	return one, nil
}
func (token *Token) GetOneByToken(tokenStr string) (*db.Token, error) {
	byToken, err := token.context.GetDb().GetTokenModel().GetOneByToken(tokenStr)
	if err != nil {
		return nil, err
	}
	token.supplement(byToken)
	return byToken, err
}

func (token *Token) supplement(st ...*db.Token) {
	mailIds := make([]uint, 0)
	stmpIds := make([]uint, 0)
	for _, d := range st {
		mailIds = append(mailIds, util.StringToUintIds(d.ReceiveEmailIds)...)
		stmpIds = append(stmpIds, d.SMTPId)
	}
	mailMap, err := token.context.GetDb().GetMailModel().GetMapByIds(mailIds)
	if err == nil {
		for _, d := range st {
			mailIds := util.StringToUintIds(d.ReceiveEmailIds)
			d.ReceiveEmails = getMails(mailIds, mailMap)
			d.ReceiveEmailsStr = getMailsStr(d.ReceiveEmails)
		}
	}

	idsMap, err := token.context.GetDb().GetSMTPModel().GetMapByIds(stmpIds)
	if err == nil {
		for _, d := range st {
			d.SMTP = idsMap[d.SMTPId]
			if d.SMTP != nil {
				d.SMTPStr = d.SMTP.Name
			}
		}
	}

}

func getMails(ids []uint, mailMap map[uint]*db.Mail) []*db.Mail {
	mails := make([]*db.Mail, 0)
	for _, id := range ids {
		v, ok := mailMap[id]
		if ok {
			mails = append(mails, v)
		}
	}
	return mails
}
func getMailsStr(mails []*db.Mail) string {
	buffer := new(bytes.Buffer)
	for _, mail := range mails {
		buffer.WriteString("," + util.FormatMail(mail.Name, mail.Mail))
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()[1:]
}
