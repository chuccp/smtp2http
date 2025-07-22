package core

import "github.com/chuccp/smtp2http/db"

type Schedule interface {
	Run(schedule *db.Schedule) error
	Stop(Id uint) error
}
