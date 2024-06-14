package api

import (
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/web"
)

type Server struct {
	context *core.Context
	core.IHttpServer
}

func NewServer() *Server {
	server := &Server{}
	httpServer := core.NewHttpServer(server.Name())
	server.IHttpServer = httpServer
	return server
}
func (s *Server) Name() string {
	return "api"
}

func (s *Server) SendMail(req *web.Request) (any, error) {

	return nil, nil
}
func (s *Server) Init(context *core.Context) {
	s.context = context
	s.IHttpServer.POST("/SendMail", s.SendMail)
}
