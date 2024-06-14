package manage

import (
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/web"
	"strconv"
)

type Stmp struct {
	context *core.Context
}

func (stmp *Stmp) getOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	one, err := stmp.context.GetDb().GetSTMPModel().GetOne(uint(atoi))
	if err != nil {
		return nil, err
	}
	return one, nil

}
func (stmp *Stmp) deleteOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	err = stmp.context.GetDb().GetSTMPModel().DeleteOne(uint(atoi))
	if err != nil {
		return nil, err
	}
	return "ok", nil

}
func (stmp *Stmp) postOne(req *web.Request) (any, error) {
	var st db.STMP
	err := req.ShouldBindBodyWithJSON(&st)
	if err != nil {
		return nil, err
	}
	err = stmp.context.GetDb().GetSTMPModel().Save(&st)
	if err != nil {
		return nil, err
	}
	return "ok", nil

}
func (stmp *Stmp) putOne(req *web.Request) (any, error) {
	var st db.STMP
	err := req.ShouldBindBodyWithJSON(&st)
	if err != nil {
		return nil, err
	}
	err = stmp.context.GetDb().GetSTMPModel().Edit(&st)
	if err != nil {
		return nil, err
	}
	return "ok", nil

}
func (stmp *Stmp) getPage(req *web.Request) (any, error) {
	page := req.GetPage()
	p, err := stmp.context.GetDb().GetSTMPModel().Page(page)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (stmp *Stmp) Init(context *core.Context, server core.IHttpServer) {
	stmp.context = context
	server.GET("/stmp/:id", stmp.getOne)
	server.DELETE("/stmp/:id", stmp.deleteOne)
	server.GET("/stmp", stmp.getPage)
	server.POST("/stmp", stmp.postOne)
	server.PUT("/stmp", stmp.putOne)
}
