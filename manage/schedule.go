package manage

import (
	"errors"
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/service"
	"github.com/chuccp/smtp2http/util"
	"github.com/chuccp/smtp2http/web"
	"strconv"
)

type Schedule struct {
	context  *core.Context
	schedule *service.Schedule
}

func (schedule *Schedule) getOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return schedule.schedule.GetOne(atoi)
}
func (schedule *Schedule) deleteOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	err = schedule.context.GetDb().GetScheduleModel().DeleteOne(uint(atoi))
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
func (schedule *Schedule) getPage(req *web.Request) (any, error) {
	page := req.GetPage()
	return schedule.schedule.GetPage(page)
}
func (schedule *Schedule) postOne(req *web.Request) (any, error) {
	var st db.Schedule
	err := req.ShouldBindBodyWithJSON(&st)
	if err != nil {
		return nil, err
	}
	err = schedule.validate(&st)
	if err != nil {
		return nil, err
	}
	err = schedule.schedule.Save(&st)
	if err != nil {
		return nil, err
	}
	err = schedule.context.GetSchedule().Run(&st)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
func (schedule *Schedule) putOne(req *web.Request) (any, error) {
	var st db.Schedule
	err := req.ShouldBindBodyWithJSON(&st)
	if err != nil {
		return nil, err
	}
	err = schedule.validate(&st)
	if err != nil {
		return nil, err
	}
	err = schedule.schedule.Edit(&st)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
func (schedule *Schedule) sendMail(req *web.Request) (any, error) {
	var st db.Schedule
	err := req.ShouldBindBodyWithJSON(&st)
	err = schedule.validate(&st)
	if err != nil {
		return nil, err
	}
	err = schedule.schedule.SendMail(&st)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
func (schedule *Schedule) validate(st *db.Schedule) error {
	err := util.ValidateURL(st.Url)
	if err != nil {
		return err
	}
	if len(st.Token) == 0 {
		return errors.New(" token  cannot be empty")
	}
	return nil
}

func (schedule *Schedule) Init(context *core.Context, server core.IHttpServer) {
	schedule.context = context
	schedule.schedule = context.GetScheduleService()
	server.GETAuth("/schedule/:id", schedule.getOne)
	server.DELETEAuth("/schedule/:id", schedule.deleteOne)
	server.GETAuth("/schedule", schedule.getPage)
	server.POSTAuth("/schedule", schedule.postOne)
	server.PUTAuth("/schedule", schedule.putOne)
	server.POSTAuth("/sendMailBySchedule", schedule.sendMail)

}
