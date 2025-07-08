package manage

import (
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/entity"
	"github.com/chuccp/smtp2http/service"
	"github.com/chuccp/smtp2http/smtp"
	"github.com/chuccp/smtp2http/util"
	"github.com/chuccp/smtp2http/web"
	"strconv"
)

type Token struct {
	context *core.Context
	token   *service.Token
	log     *service.Log
}

func (token *Token) getOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return token.token.GetOne(atoi)
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
	return token.token.GetPage(page)
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
	st.ReceiveEmailIds = util.DeduplicateIds(st.ReceiveEmailIds)
	err = token.context.GetDb().GetTokenModel().Edit(&st)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}

func (token *Token) sendMail(req *web.Request) (any, error) {
	var sendMail entity.SendMail
	err := req.ShouldBindBodyWithJSON(&sendMail)
	if err != nil {
		return nil, err
	}
	byToken, err := token.token.GetOne(int(sendMail.SMTPId))
	if err != nil {
		return nil, err
	}
	if len(sendMail.Subject) == 0 {
		sendMail.Subject = byToken.Subject
	}
	err = smtp.SendContentMsgByRecipients(byToken.SMTP, sendMail.Recipients, sendMail.Subject, sendMail.Content)
	if err != nil {
		token.log.ContentError(byToken.SMTP, byToken.ReceiveEmails, byToken.Token, sendMail.Subject, sendMail.Content, err)
		return nil, err
	} else {
		token.log.ContentSuccess(byToken.SMTP, byToken.ReceiveEmails, byToken.Token, sendMail.Subject, sendMail.Content)
	}
	return "ok", nil
}

func (token *Token) Init(context *core.Context, server core.IHttpServer) {
	token.context = context
	token.token = context.GetTokenService()
	token.log = context.GetLogService()
	server.GETAuth("/token/:id", token.getOne)
	server.DELETEAuth("/token/:id", token.deleteOne)
	server.GETAuth("/token", token.getPage)
	server.POSTAuth("/token", token.postOne)
	server.PUTAuth("/token", token.putOne)
	server.POSTAuth("/sendMailByToken", token.sendMail)

}
