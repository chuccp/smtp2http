package api

import (
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/service"
	"github.com/chuccp/smtp2http/web"
)

type Server struct {
	context *core.Context
	core.IHttpServer
	token *service.Token
	log   *service.Log
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
func (s *Server) Start() {}

func (s *Server) SendMail(req *web.Request) (any, error) {
	return s.token.SendMailByToken(req)
}

func (s *Server) Init(context *core.Context) {
	s.context = context
	s.token = context.GetTokenService()
	s.log = context.GetLogService()
	s.POST("/sendMail", s.SendMail)
	s.GET("/sendMail", s.SendMail)
}
