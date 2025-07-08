package service

import (
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/util"
	"github.com/chuccp/smtp2http/web"
)

type Token struct {
	db *db.DB
}

func NewToken(db *db.DB) *Token {
	return &Token{db: db}
}
func (token *Token) GetPage(page *web.Page) (any, error) {
	p, err := token.db.GetTokenModel().Page(page)
	if err != nil {
		return nil, err
	}
	token.supplementToken(p.List...)
	return p, nil
}

func (token *Token) supplementToken(st ...*db.Token) {
	mailIds := make([]uint, 0)
	stmpIds := make([]uint, 0)
	for _, d := range st {
		d.Name = d.Token
		mailIds = append(mailIds, util.StringToUintIds(d.ReceiveEmailIds)...)
		stmpIds = append(stmpIds, d.SMTPId)
	}
	mailMap, err := token.db.GetMailModel().GetMapByIds(mailIds)
	if err == nil {
		for _, d := range st {
			mailIds := util.StringToUintIds(d.ReceiveEmailIds)
			d.ReceiveEmails = db.GetMails(mailIds, mailMap)
			d.ReceiveEmailsStr = db.GetMailsStr(d.ReceiveEmails)
		}
	}
	idsMap, err := token.db.GetSMTPModel().GetMapByIds(stmpIds)
	if err == nil {
		for _, d := range st {
			d.SMTP = idsMap[d.SMTPId]
			if d.SMTP != nil {
				d.SMTPStr = d.SMTP.Name
			}
		}
	}
}

func (token *Token) GetOne(id int) (*db.Token, error) {
	one, err := token.db.GetTokenModel().GetOne(uint(id))
	if err != nil {
		return nil, err
	}
	token.supplementToken(one)
	return one, nil
}
func (token *Token) GetOneByToken(tokenStr string) (*db.Token, error) {
	byToken, err := token.db.GetTokenModel().GetOneByToken(tokenStr)
	if err != nil {
		return nil, err
	}
	token.supplementToken(byToken)
	return byToken, err
}
