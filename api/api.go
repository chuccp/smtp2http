package api

import (
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/service"
	"github.com/chuccp/d-mail/stmp"
	"github.com/chuccp/d-mail/web"
)

type Server struct {
	context *core.Context
	core.IHttpServer
	token *service.Token
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
	token := req.FormValue("token")
	content := req.FormValue("content")
	subject := req.FormValue("subject")
	byToken, err := s.token.GetOneByToken(token)
	if err != nil {
		return nil, err
	}
	if req.IsMultipartForm() {
		form, err := req.MultipartForm()
		if err != nil {
			return nil, err
		}
		fileHeaders, ok := form.File["files"]
		if ok {
			for _, fileHeader := range fileHeaders {
				web.SaveUploadedFile(fileHeader, fileHeader.Filename)
			}
		}
	} else {
		if len(subject) == 0 {
			subject = byToken.Subject
		}
		err := stmp.SendContentMsg(byToken.STMP, byToken.ReceiveEmails, subject, content)
		if err != nil {
			return nil, err
		}
	}
	return "ok", nil
}
func (s *Server) Init(context *core.Context) {
	s.context = context
	s.token = service.NewToken(context)
	s.IHttpServer.POST("/sendMail", s.SendMail)
	s.IHttpServer.GET("/sendMail", s.SendMail)
}
