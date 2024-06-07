package api

import (
	"bytes"
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/util"
	"github.com/chuccp/d-mail/web"
	"strconv"
)

type Token struct {
	context *core.Context
}

func (token *Token) getOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	one, err := token.context.GetDb().GetTokenModel().GetOne(uint(atoi))
	if err != nil {
		return nil, err
	}
	token.supplement(one)
	return one, nil
}

func (token *Token) deleteOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	err = token.context.GetDb().GetTokenModel().DeleteOne(uint(atoi))
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
func (token *Token) getPage(req *web.Request) (any, error) {
	page := req.GetPage()
	p, err := token.context.GetDb().GetTokenModel().Page(page)
	if err != nil {
		return nil, err
	}
	token.supplement(p.List...)
	return p, nil
}

func (token *Token) supplement(st ...*db.Token) {
	mailIds := make([]uint, 0)
	stmpIds := make([]uint, 0)
	for _, d := range st {
		mailIds = append(mailIds, util.StringToUintIds(d.ReceiveEmailIds)...)
		stmpIds = append(stmpIds, d.STMPId)
	}
	mailMap, err := token.context.GetDb().GetMailModel().GetMapByIds(mailIds)
	if err == nil {
		for _, d := range st {
			mailIds := util.StringToUintIds(d.ReceiveEmailIds)
			d.ReceiveEmails = getMails(mailIds, mailMap)
			d.ReceiveEmailsStr = getMailsStr(d.ReceiveEmails)
		}
	}

	idsMap, err := token.context.GetDb().GetSTMPModel().GetMapByIds(stmpIds)
	if err == nil {
		for _, d := range st {
			d.STMP = idsMap[d.STMPId]
			if d.STMP != nil {
				d.STMPStr = d.STMP.Name
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
		buffer.WriteString(";" + mail.Name + ":[" + mail.Mail + "]")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()[1:]
}

func (token *Token) postOne(req *web.Request) (any, error) {
	var st db.Token
	err := req.ShouldBindBodyWithJSON(&st)
	if err != nil {
		return nil, err
	}
	err = token.context.GetDb().GetTokenModel().Save(&st)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
func (token *Token) putOne(req *web.Request) (any, error) {
	var st db.Token
	err := req.ShouldBindBodyWithJSON(&st)
	if err != nil {
		return nil, err
	}
	err = token.context.GetDb().GetTokenModel().Edit(&st)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
func (token *Token) Init(context *core.Context) {
	token.context = context
	context.GET("/token/:id", token.getOne)
	context.DELETE("/token/:id", token.deleteOne)
	context.GET("/token", token.getPage)
	context.POST("/token", token.postOne)
	context.PUT("/token", token.putOne)
}
