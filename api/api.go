package api

import (
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/service"
	"github.com/chuccp/d-mail/stmp"
	"github.com/chuccp/d-mail/util"
	"github.com/chuccp/d-mail/web"
	"os"
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

func (s *Server) SendMail(req *web.Request) (any, error) {
	token := req.FormValue("token")
	content := req.FormValue("content")
	subject := req.FormValue("subject")
	byToken, err := s.token.GetOneByToken(token)
	if err != nil {
		return nil, err
	}
	if req.IsMultipartForm() {
		chachePath := s.context.GetConfig().GetStringOrDefault("core", "cache-path", "cache")
		form, err := req.MultipartForm()
		if err != nil {
			return nil, err
		}
		fileHeaders, ok := form.File["file"]
		if ok {
			files := make([]*stmp.File, 0)
			for _, fileHeader := range fileHeaders {
				filePath := util.GetCachePath(chachePath, fileHeader.Filename)
				err := web.SaveUploadedFile(fileHeader, filePath)
				if err != nil {
					return nil, err
				}
				file, err := os.Open(filePath)
				if err != nil {
					return nil, err
				}
				files = append(files, &stmp.File{File: file, Name: fileHeader.Filename})
			}
			if len(files) > 0 {
				err := stmp.SendFilesMsg(byToken.STMP, byToken.ReceiveEmails, files, subject, content)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		if len(subject) == 0 {
			subject = byToken.Subject
		}
		err := stmp.SendContentMsg(byToken.STMP, byToken.ReceiveEmails, subject, content)
		if err != nil {
			s.log.ContentError(byToken.STMP, byToken.ReceiveEmails, subject, content, err)
			return nil, err
		} else {
			s.log.ContentSuccess(byToken.STMP, byToken.ReceiveEmails, subject, content)
		}
	}
	return "ok", nil
}

func (s *Server) Init(context *core.Context) {
	s.context = context
	s.token = service.NewToken(context)
	s.log = service.NewLog(context)
	s.IHttpServer.POST("/sendMail", s.SendMail)
	s.IHttpServer.GET("/sendMail", s.SendMail)
}
