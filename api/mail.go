package api

import (
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/web"
	"strconv"
)

type Mail struct {
	context *core.Context
}

func (mail *Mail) getOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	one, err := mail.context.GetDb().GetMailModel().GetOne(uint(atoi))
	if err != nil {
		return nil, err
	}
	return one, nil
}

func (mail *Mail) deleteOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	err = mail.context.GetDb().GetMailModel().DeleteOne(uint(atoi))
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
func (mail *Mail) getPage(req *web.Request) (any, error) {
	page := req.GetPage()
	p, err := mail.context.GetDb().GetMailModel().Page(page)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (mail *Mail) postOne(req *web.Request) (any, error) {
	var st db.Mail
	err := req.ShouldBindBodyWithJSON(&st)
	if err != nil {
		return nil, err
	}
	err = mail.context.GetDb().GetMailModel().Save(&st)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
func (mail *Mail) putOne(req *web.Request) (any, error) {
	var st db.Mail
	err := req.ShouldBindBodyWithJSON(&st)
	if err != nil {
		return nil, err
	}
	err = mail.context.GetDb().GetMailModel().Edit(&st)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
func (mail *Mail) Init(context *core.Context) {
	mail.context = context
	context.GET("/mail/:id", mail.getOne)
	context.DELETE("/mail/:id", mail.deleteOne)
	context.GET("/mail", mail.getPage)
	context.POST("/mail", mail.postOne)
	context.PUT("/mail", mail.putOne)
}
