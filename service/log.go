package service

import (
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/util"
	"go.uber.org/zap/buffer"
)

type Log struct {
	context *core.Context
}

func NewLog(context *core.Context) *Log {
	return &Log{context: context}
}

func (a *Log) ContentSuccess(stmp *db.STMP, mails []*db.Mail, subject, bodyString string) error {
	return a.log(stmp, mails, subject, bodyString, db.SUCCESS, nil)
}
func (a *Log) log(stmp *db.STMP, mails []*db.Mail, subject, bodyString string, status byte, err error) error {
	var lg db.Log
	lg.STMP = util.FormatMail(stmp.Username, stmp.Mail)
	buffer := new(buffer.Buffer)
	for _, mail := range mails {
		buffer.AppendString(",")
		buffer.AppendString(util.FormatMail(mail.Name, mail.Mail))
	}
	if buffer.Len() > 0 {
		lg.Mail = buffer.String()[1:]
	}
	lg.Subject = subject
	lg.Content = bodyString
	if status == db.SUCCESS {
		lg.Result = "success"
		lg.Status = status
	}
	if status == db.ERROR {
		lg.Result = err.Error()
		lg.Status = status
	}
	return a.context.GetDb().GetLogModel().Save(&lg)
}
func (a *Log) ContentError(stmp *db.STMP, mails []*db.Mail, subject, bodyString string, err error) error {
	return a.log(stmp, mails, subject, bodyString, db.ERROR, err)
}
