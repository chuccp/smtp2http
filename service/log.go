package service

import (
	"errors"
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/stmp"
	"github.com/chuccp/d-mail/util"
	"go.uber.org/zap/buffer"
)

type Log struct {
	context *core.Context
}

func NewLog(context *core.Context) *Log {
	return &Log{context: context}
}

func (a *Log) ContentSuccess(stmp *db.STMP, mails []*db.Mail, token string, subject, bodyString string) error {
	return a.log(stmp, mails, token, subject, bodyString, nil, db.SUCCESS, nil)
}
func (a *Log) log(st *db.STMP, mails []*db.Mail, token string, subject, bodyString string, files []*stmp.File, status byte, err error) error {
	var lg db.Log
	lg.Token = token
	lg.STMP = util.FormatMail(st.Username, st.Mail)
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
	} else {
		if files != nil && len(files) > 0 {

		}
		var ee *stmp.UserNotFoundError
		ok := errors.As(err, &ee)
		if ok {
			lg.Status = db.WARM
			lg.Result = err.Error()
		} else {
			lg.Result = err.Error()
			lg.Status = status
		}
	}
	return a.context.GetDb().GetLogModel().Save(&lg)
}
func (a *Log) FilesError(stmp *db.STMP, mails []*db.Mail, files []*stmp.File, token string, subject, bodyString string, err error) error {
	return a.log(stmp, mails, token, subject, bodyString, files, db.ERROR, err)
}
func (a *Log) FilesSuccess(stmp *db.STMP, mails []*db.Mail, files []*stmp.File, token string, subject, bodyString string) error {
	return a.log(stmp, mails, token, subject, bodyString, files, db.SUCCESS, nil)
}

func (a *Log) ContentError(stmp *db.STMP, mails []*db.Mail, token string, subject string, bodyString string, err error) error {
	return a.log(stmp, mails, token, subject, bodyString, nil, db.ERROR, err)
}
