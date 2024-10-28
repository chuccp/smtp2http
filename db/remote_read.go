package db

import "time"

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
