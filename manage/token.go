package manage

import (
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/service"
	"github.com/chuccp/d-mail/util"
	"github.com/chuccp/d-mail/web"
	"strconv"
)

type Token struct {
	context *core.Context
	token   *service.Token
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
func (token *Token) Init(context *core.Context, server core.IHttpServer) {
	token.context = context
	token.token = service.NewToken(context)
	server.GETAuth("/token/:id", token.getOne)
	server.DELETEAuth("/token/:id", token.deleteOne)
	server.GETAuth("/token", token.getPage)
	server.POSTAuth("/token", token.postOne)
	server.PUTAuth("/token", token.putOne)
}
