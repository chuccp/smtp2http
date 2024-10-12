package api

import (
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/entity"
	"github.com/chuccp/smtp2http/service"
	"github.com/chuccp/smtp2http/smtp"
	"github.com/chuccp/smtp2http/util"
	"github.com/chuccp/smtp2http/web"
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
	var sendMailApi entity.SendMailApi
	if util.ContainsAnyIgnoreCase(req.GetContext().ContentType(), "application/json") {
		err := req.ShouldBindBodyWithJSON(&sendMailApi)
		if err != nil {
			return nil, err
		}
	} else {
		sendMailApi.Token = req.FormValue("token")
		sendMailApi.Content = req.FormValue("content")
		sendMailApi.Subject = req.FormValue("subject")
		sendMailApi.Recipients = util.SplitAndDeduplicate(req.FormValue("recipients"), ",")
	}
	byToken, err := s.token.GetOneByToken(sendMailApi.Token)
	if err != nil {
		return nil, err
	}
	for _, mail := range sendMailApi.Recipients {
		byToken.ReceiveEmails = append(byToken.ReceiveEmails, &db.Mail{Mail: mail})
	}
	if req.IsMultipartForm() {
		cachePath := s.context.GetConfig().GetStringOrDefault("core", "cachePath", ".cache")
		form, err := req.MultipartForm()
		if err != nil {
			return nil, err
		}
		fileHeaders, ok := form.File["files"]
		if ok {
			files := make([]*smtp.File, 0)
			for _, fileHeader := range fileHeaders {
				filePath := util.GetCachePath(cachePath, fileHeader.Filename)
				err := web.SaveUploadedFile(fileHeader, filePath)
				if err != nil {
					return nil, err
				}
				file, err := os.Open(filePath)
				if err != nil {
					return nil, err
				}
				files = append(files, &smtp.File{File: file, Name: fileHeader.Filename, FilePath: filePath})
			}
			if len(files) > 0 {
				err := smtp.SendFilesMsg(byToken.SMTP, byToken.ReceiveEmails, files, sendMailApi.Subject, sendMailApi.Content)
				if err != nil {
					s.log.FilesError(byToken.SMTP, byToken.ReceiveEmails, files, sendMailApi.Token, sendMailApi.Subject, sendMailApi.Content, err)
					return nil, err
				} else {
					s.log.FilesSuccess(byToken.SMTP, byToken.ReceiveEmails, files, sendMailApi.Token, sendMailApi.Subject, sendMailApi.Content)
				}
			}
		}
	} else {
		if len(sendMailApi.Subject) == 0 {
			sendMailApi.Subject = byToken.Subject
		}
		err := smtp.SendContentMsg(byToken.SMTP, byToken.ReceiveEmails, sendMailApi.Subject, sendMailApi.Content)
		if err != nil {
			s.log.ContentError(byToken.SMTP, byToken.ReceiveEmails, sendMailApi.Token, sendMailApi.Subject, sendMailApi.Content, err)
			return nil, err
		} else {
			s.log.ContentSuccess(byToken.SMTP, byToken.ReceiveEmails, sendMailApi.Token, sendMailApi.Subject, sendMailApi.Content)
		}
	}
	return "ok", nil
}

func (s *Server) Init(context *core.Context) {
	s.context = context
	s.token = service.NewToken(context)
	s.log = service.NewLog(context)
	s.POST("/sendMail", s.SendMail)
	s.GET("/sendMail", s.SendMail)
}
