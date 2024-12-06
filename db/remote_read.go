package db

import (
	"github.com/chuccp/smtp2http/web"
	"gorm.io/gorm"
	"time"
)

type RemoteRead struct {
	Id          uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Url         string    `gorm:"column:url" json:"url"`
	Header      string    `gorm:"column:Header" json:"header"`
	Method      string    `gorm:"column:method" json:"method"`
	Body        string    `gorm:"column:body" json:"body"`
	BodyFormat  string    `gorm:"column:body_format" json:"bodyFormat"`
	ContentType string    `gorm:"column:content_type" json:"contentType"`
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (token *RemoteRead) SetCreateTime(createTime time.Time) {
	token.CreateTime = createTime
}
func (token *RemoteRead) SetUpdateTime(updateTIme time.Time) {
	token.UpdateTime = updateTIme
}
func (token *RemoteRead) GetId() uint {
	return token.Id
}
func (token *RemoteRead) SetId(id uint) {
	token.Id = id
}

type RemoteReadModel struct {
	*Model[*RemoteRead]
	db        *gorm.DB
	tableName string
}

func NewRemoteReadModel(db *gorm.DB, tableName string) *RemoteReadModel {
	return &RemoteReadModel{db: db, tableName: tableName, Model: NewModel[*RemoteRead](db, tableName)}
}
func (a *RemoteReadModel) CreateTable() error {
	return a.Model.CreateTable(&RemoteRead{})
}
func (a *RemoteReadModel) DeleteTable() error {
	return a.Model.DeleteTable(&RemoteRead{})
}

func (a *RemoteReadModel) Save(remoteRead *RemoteRead) error {
	return a.Model.Save(remoteRead)
}
func (a *RemoteReadModel) GetOne(id uint) (*RemoteRead, error) {
	var remoteRead RemoteRead
	err := a.Model.GetOne(id, &remoteRead)
	if err != nil {
		return nil, err
	}
	return &remoteRead, nil
}

func (a *RemoteReadModel) DeleteOne(id uint) error {
	return a.Model.DeleteOne(id, &RemoteRead{})
}

func (a *RemoteReadModel) NewModel(db *gorm.DB) *RemoteReadModel {
	return &RemoteReadModel{db: db, tableName: a.tableName}
}
func (a *RemoteReadModel) Page(page *web.Page) (*Page[*RemoteRead], error) {
	var remoteReads []*RemoteRead
	num, err := a.Model.Page(page, &remoteReads)
	if err != nil {
		return nil, err
	}
	return ToPage[*RemoteRead](num, remoteReads), nil
}
