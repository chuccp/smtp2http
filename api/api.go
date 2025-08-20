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
func (s *Server) Start() {}

func (s *Server) SendMail(req *web.Request) (any, error) {

	var sendMailApi entity.SendMailApi
	var byToken *db.Token
	files := make([]*smtp.File, 0)
	err := func() (err error) {
		if util.ContainsAnyIgnoreCase(req.GetContext().ContentType(), "application/json") {
			err = req.ShouldBindBodyWithJSON(&sendMailApi)
			if err != nil {
				return err
			}
		} else {
			sendMailApi.Token = req.FormValue("token")
			sendMailApi.Content = req.FormValue("content")
			sendMailApi.Subject = req.FormValue("subject")
			sendMailApi.Recipients = util.SplitAndDeduplicate(req.FormValue("recipients"), ",")
		}
		byToken, err = s.token.GetOneByToken(sendMailApi.Token)
		if err != nil {
			return err
		}
		for _, mail := range sendMailApi.Recipients {
			byToken.ReceiveEmails = append(byToken.ReceiveEmails, &db.Mail{Mail: mail})
		}
		if len(sendMailApi.Subject) == 0 {
			sendMailApi.Subject = byToken.Subject
		}

		if req.IsMultipartForm() || len(sendMailApi.Files) > 0 {
			cachePath := s.context.GetConfig().GetStringOrDefault("core", "cachePath", ".cache")
			form, err := req.MultipartForm()
			if err != nil {
				return err
			}
			fileHeaders, ok := form.File["files"]
			if ok {
				for _, fileHeader := range fileHeaders {
					filePath := util.GetCachePath(cachePath, fileHeader.Filename)
					err := web.SaveUploadedFile(fileHeader, filePath)
					if err != nil {
						return err
					}
					file, err := os.Open(filePath)
					if err != nil {
						return err
					}
					files = append(files, &smtp.File{File: file, Name: fileHeader.Filename, FilePath: filePath})
				}
			}
			for _, file := range sendMailApi.Files {
				if len(file.Data) == 0 {
					continue
				}
				base64, err := util.DecodeBase64(file.Data)
				if err != nil {
					return err
				}
				if len(file.Name) == 0 {
					file.Name, err = util.CalculateMD5(base64)
					if err != nil {
						return err
					}
				}
				filePath := util.GetCachePath(cachePath, file.Name)
				err = util.WriteFile(base64, filePath)
				if err != nil {
					return err
				}
				file, err := os.Open(filePath)
				if err != nil {
					return err
				}
				files = append(files, &smtp.File{File: file, Name: file.Name(), FilePath: filePath})
			}
		}
		return nil
	}()
	if err == nil {
		err = smtp.SendAllMsg(byToken.SMTP, byToken.ReceiveEmails, files, sendMailApi.Subject, sendMailApi.Content)
	}
	err = s.log.Log(byToken.SMTP, byToken.ReceiveEmails, files, sendMailApi.Token, sendMailApi.Subject, sendMailApi.Content, err)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}

func (s *Server) Init(context *core.Context) {
	s.context = context
	s.token = context.GetTokenService()
	s.log = context.GetLogService()
	s.POST("/sendMail", s.SendMail)
	s.GET("/sendMail", s.SendMail)
}
