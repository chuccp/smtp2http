package db

import (
	"github.com/chuccp/d-mail/web"
	"gorm.io/gorm"
	"time"
)

type STMP struct {
	Id         uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Host       string    `gorm:"column:host" json:"host"`
	Port       int       `gorm:"column:port" json:"port"`
	Mail       string    `gorm:"column:mail" json:"mail"`
	Username   string    `gorm:"column:username" json:"username"`
	Password   string    `gorm:"column:password"  json:"password"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (stmp *STMP) SetCreateTime(createTime time.Time) {
	stmp.CreateTime = createTime
}
func (stmp *STMP) SetUpdateTime(updateTime time.Time) {
	stmp.UpdateTime = updateTime
}
func (stmp *STMP) GetId() uint {
	return stmp.Id
}
func (stmp *STMP) SetId(id uint) {
	stmp.Id = id
}

type STMPModel struct {
	*Model[*STMP]
	db        *gorm.DB
	tableName string
}

func NewSTMPModel(db *gorm.DB, tableName string) *STMPModel {
	return &STMPModel{db: db, tableName: tableName, Model: NewModel[*STMP](db, tableName)}
}

func (a *STMPModel) CreateTable() error {
	return a.Model.CreateTable(&STMP{})
}
func (a *STMPModel) DeleteTable() error {
	return a.Model.DeleteTable(&STMP{})
}

func (a *STMPModel) Save(stmp *STMP) error {
	return a.Model.Save(stmp)
}
func (a *STMPModel) GetOne(id uint) (*STMP, error) {
	var stmp STMP
	err := a.Model.GetOne(id, &stmp)
	if err != nil {
		return nil, err
	}
	return &stmp, nil
}

func (a *STMPModel) DeleteOne(id uint) error {
	return a.Model.DeleteOne(id, &STMP{})
}

func (a *STMPModel) Edit(stmp *STMP) error {
	return a.Model.Edit(stmp)
}

func (a *STMPModel) NewModel(db *gorm.DB) *STMPModel {
	return &STMPModel{db: db, tableName: a.tableName}
}
func (a *STMPModel) Page(page *web.Page) (*Page[*STMP], error) {
	var stmps []*STMP
	num, err := a.Model.Page(page, &stmps)
	if err != nil {
		return nil, err
	}
	return ToPage[*STMP](num, stmps), nil
}
