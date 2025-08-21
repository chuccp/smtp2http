package service

import (
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/entity"
	"github.com/chuccp/smtp2http/smtp"
	"github.com/chuccp/smtp2http/util"
	"github.com/chuccp/smtp2http/web"
	"os"
)

type Token struct {
	db        *db.DB
	log       *Log
	cachePath string
}

func NewToken(db *db.DB, cachePath string) *Token {
	return &Token{db: db, log: NewLog(db), cachePath: cachePath}
}
func (token *Token) GetPage(page *web.Page) (any, error) {
	p, err := token.db.GetTokenModel().Page(page)
	if err != nil {
		return nil, err
	}
	token.supplementToken(p.List...)
	return p, nil
}

func (token *Token) supplementToken(st ...*db.Token) {
	mailIds := make([]uint, 0)
	stmpIds := make([]uint, 0)
	for _, d := range st {
		d.Name = d.Token
		mailIds = append(mailIds, util.StringToUintIds(d.ReceiveEmailIds)...)
		stmpIds = append(stmpIds, d.SMTPId)
	}
	mailMap, err := token.db.GetMailModel().GetMapByIds(mailIds)
	if err == nil {
		for _, d := range st {
			mailIds := util.StringToUintIds(d.ReceiveEmailIds)
			d.ReceiveEmails = db.GetMails(mailIds, mailMap)
			d.ReceiveEmailsStr = db.GetMailsStr(d.ReceiveEmails)
		}
	}
	idsMap, err := token.db.GetSMTPModel().GetMapByIds(stmpIds)
	if err == nil {
		for _, d := range st {
			d.SMTP = idsMap[d.SMTPId]
			if d.SMTP != nil {
				d.SMTPStr = d.SMTP.Name
			}
		}
	}
}

func (token *Token) SendMailByToken(req *web.Request) (any, error) {
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
		byToken, err = token.GetOneByToken(sendMailApi.Token)
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
			form, err := req.MultipartForm()
			if err != nil {
				return err
			}
			fileHeaders, ok := form.File["files"]
			if ok {
				for _, fileHeader := range fileHeaders {
					filePath := util.GetCachePath(token.cachePath, fileHeader.Filename)
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
				filePath := util.GetCachePath(token.cachePath, file.Name)
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
	err = token.log.Log(byToken.SMTP, byToken.ReceiveEmails, files, sendMailApi.Token, sendMailApi.Subject, sendMailApi.Content, err)
	if err != nil {
		return nil, err
	}
	return "ok", nil

}

func (token *Token) GetOne(id int) (*db.Token, error) {
	one, err := token.db.GetTokenModel().GetOne(uint(id))
	if err != nil {
		return nil, err
	}
	token.supplementToken(one)
	return one, nil
}
func (token *Token) GetOneByToken(tokenStr string) (*db.Token, error) {
	byToken, err := token.db.GetTokenModel().GetOneByToken(tokenStr)
	if err != nil {
		return nil, err
	}
	token.supplementToken(byToken)
	return byToken, err
}
