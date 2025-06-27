package db

import (
	"github.com/chuccp/smtp2http/web"
	"gorm.io/gorm"
	"time"
)

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Schedule struct {
	Id                uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name              string    `gorm:"column:name" json:"name"`
	Token             string    `gorm:"column:token" json:"token"`
	Cron              string    `gorm:"column:cron" json:"cron"`
	Url               string    `gorm:"column:url" json:"url"`
	Method            string    `gorm:"column:method" json:"method"`
	HeaderStr         string    `gorm:"column:header_str" json:"headerStr"`
	Headers           []*Header `gorm:"-" json:"headers"`
	Body              string    `gorm:"column:body" json:"body"`
	UseTemplate       bool      `gorm:"column:use_template" json:"useTemplate"`
	Template          string    `gorm:"column:template" json:"template"`
	IsUse             bool      `gorm:"column:is_use" json:"isUse"`
	IsOnlySendByError bool      `gorm:"column:is_only_send_by_error" json:"isOnlySendByError"`
	CreateTime        time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime        time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (schedule *Schedule) SetCreateTime(createTime time.Time) {
	schedule.CreateTime = createTime
}
func (schedule *Schedule) SetUpdateTime(updateTIme time.Time) {
	schedule.UpdateTime = updateTIme
}
func (schedule *Schedule) GetId() uint {
	return schedule.Id
}
func (schedule *Schedule) SetId(id uint) {
	schedule.Id = id
}

type ScheduleModel struct {
	*Model[*Schedule]
	db        *gorm.DB
	tableName string
}

func NewScheduleModel(db *gorm.DB, tableName string) *ScheduleModel {
	return &ScheduleModel{db: db, tableName: tableName, Model: NewModel[*Schedule](db, tableName)}
}
func (a *ScheduleModel) CreateTable() error {
	return a.Model.CreateTable(&Schedule{})
}
func (a *ScheduleModel) DeleteTable() error {
	return a.Model.DeleteTable(&Schedule{})
}

func (a *ScheduleModel) Save(remoteRead *Schedule) error {
	return a.Model.Save(remoteRead)
}
func (a *ScheduleModel) GetOne(id uint) (*Schedule, error) {
	var schedule Schedule
	err := a.Model.GetOne(id, &schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}
func (a *ScheduleModel) FindAllByUse() ([]*Schedule, error) {
	var schedules []*Schedule
	tx := a.db.Table(a.tableName).Where("`is_use`=1 ").Find(&schedules)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return schedules, nil
}

func (a *ScheduleModel) DeleteOne(id uint) error {
	return a.Model.DeleteOne(id, &Schedule{})
}

func (a *ScheduleModel) NewModel(db *gorm.DB) *ScheduleModel {
	return &ScheduleModel{db: db, tableName: a.tableName}
}
func (a *ScheduleModel) Page(page *web.Page) (*Page[*Schedule], error) {
	var schedules []*Schedule
	num, err := a.Model.Page(page, &schedules)
	if err != nil {
		return nil, err
	}
	return ToPage[*Schedule](num, schedules), nil
}
