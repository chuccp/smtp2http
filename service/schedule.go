package service

import (
	"errors"
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/smtp"
	"github.com/chuccp/smtp2http/web"
)

type Schedule struct {
	db    *db.DB
	token *Token
}

func NewSchedule(db *db.DB, token *Token) *Schedule {
	return &Schedule{db: db, token: token}
}
func (schedule *Schedule) GetPage(page *web.Page) (any, error) {
	return schedule.db.GetScheduleModel().Page(page)
}
func (schedule *Schedule) Edit(sd *db.Schedule) error {
	v, err := schedule.db.GetTokenModel().GetOneByToken(sd.Token)
	if err != nil {
		return err
	}
	if v == nil {
		return errors.New("token not found")
	}
	return schedule.db.GetScheduleModel().Edit(sd)
}
func (schedule *Schedule) Save(sd *db.Schedule) error {
	v, err := schedule.db.GetTokenModel().GetOneByToken(sd.Token)
	if err != nil {
		return err
	}
	if v == nil {
		return errors.New("token not found")
	}
	return schedule.db.GetScheduleModel().Save(sd)

}

func (schedule *Schedule) GetOne(id int) (*db.Schedule, error) {
	return schedule.db.GetScheduleModel().GetOne(uint(id))
}

func (schedule *Schedule) SendMail(sd *db.Schedule) error {
	byToken, err := schedule.token.GetOneByToken(sd.Token)
	if err != nil {
		return err
	}
	if byToken == nil {
		return errors.New("token not found")
	}
	_, err = smtp.SendAPIMail(sd, byToken.SMTP, byToken.ReceiveEmails)
	return err
}
