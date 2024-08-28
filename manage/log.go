package manage

import (
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/web"
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
	one, err := log.context.GetDb().GetLogModel().GetOne(uint(atoi))
	if err != nil {
		return nil, err
	}
	return one, nil
}

func (log *Log) getPage(req *web.Request) (any, error) {
	page := req.GetPage()
	p, err := log.context.GetDb().GetLogModel().Page(page)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (log *Log) downLoad(req *web.Request) (any, error) {
	return nil, nil
}
func (log *Log) Init(context *core.Context, server core.IHttpServer) {
	log.context = context
	server.GETAuth("/log/:id", log.getOne)
	server.GETAuth("/log", log.getPage)
	server.GETAuth("/downLoad", log.downLoad)
}
