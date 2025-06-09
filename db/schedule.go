package db

import (
	"github.com/chuccp/smtp2http/web"
	"gorm.io/gorm"
	"time"
)

type Schedule struct {
	Id          uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Url         string    `gorm:"column:url" json:"url"`
	Header      string    `gorm:"column:Header" json:"header"`
	Method      string    `gorm:"column:method" json:"method"`
	ContentType string    `gorm:"column:content_type" json:"contentType"`
	Body        string    `gorm:"column:body" json:"body"`
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`
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
