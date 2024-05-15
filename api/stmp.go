package api

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
func (stmp *Stmp) Init(context *core.Context) {
	stmp.context = context
	context.GET("/stmp/:id", stmp.getOne)
	context.DELETE("/stmp/:id", stmp.deleteOne)
	context.GET("/stmp", stmp.getPage)
	context.POST("/stmp", stmp.postOne)
	context.PUT("/stmp", stmp.putOne)
}
