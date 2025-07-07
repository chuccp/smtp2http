package service

import (
	"errors"
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/smtp"
	"github.com/chuccp/smtp2http/web"
)

type Schedule struct {
	context *core.Context
	token   *Token
}

func NewSchedule(context *core.Context) *Schedule {
	return &Schedule{context: context, token: NewToken(context)}
}
func (schedule *Schedule) GetPage(page *web.Page) (any, error) {
	return schedule.context.GetDb().GetScheduleModel().Page(page)
}
func (schedule *Schedule) Edit(sd *db.Schedule) error {
	token := sd.Token
	v, err := schedule.token.QueryOneByToken(token)
	if v == nil {
		return errors.New("token not found")
	}
	if err != nil {
		return err
	}
	return schedule.context.GetDb().GetScheduleModel().Edit(sd)
}
func (schedule *Schedule) Save(sd *db.Schedule) error {
	token := sd.Token
	v, err := schedule.token.QueryOneByToken(token)
	if v == nil {
		return errors.New("token not found")
	}
	if err != nil {
		return err
	}
	return schedule.context.GetDb().GetScheduleModel().Save(sd)

}

func (schedule *Schedule) GetOne(id int) (*db.Schedule, error) {
	return schedule.context.GetDb().GetScheduleModel().GetOne(uint(id))
}

func (schedule *Schedule) SendMail(sd *db.Schedule) error {
	token := sd.Token
	byToken, err := schedule.token.GetOneByToken(token)
	if byToken == nil {
		return errors.New("token not found")
	}
	if err != nil {
		return err
	}
	return smtp.SendAPIMail(sd, byToken.SMTP, byToken.ReceiveEmails)
}
