package manage

import (
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/web"
	"strconv"
)

type Log struct {
	context *core.Context
}

func (log *Log) getOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	one, err := log.context.GetDb().GetMailModel().GetOne(uint(atoi))
	if err != nil {
		return nil, err
	}
	return one, nil
}

func (log *Log) deleteOne(req *web.Request) (any, error) {
	return nil, nil
}

func (log *Log) getPage(req *web.Request) (any, error) {
	page := req.GetPage()
	p, err := log.context.GetDb().GetLogModel().Page(page)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (log *Log) Init(context *core.Context, server core.IHttpServer) {
	log.context = context
	server.GET("/log/:id", log.getOne)
	server.DELETE("/log/:id", log.deleteOne)
	server.GET("/log", log.getPage)
}
