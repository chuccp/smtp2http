package service

import (
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/web"
)

type Schedule struct {
	context *core.Context
}

func NewSchedule(context *core.Context) *Schedule {
	return &Schedule{context: context}
}
func (schedule *Schedule) GetPage(page *web.Page) (any, error) {
	return schedule.context.GetDb().GetScheduleModel().Page(page)
}
func (schedule *Schedule) GetOne(id int) (*db.Schedule, error) {
	return schedule.context.GetDb().GetScheduleModel().GetOne(uint(id))
}
