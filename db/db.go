package db

import (
	"fmt"
	"github.com/chuccp/d-mail/util"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var NoDatabase = &NoDatabaseError{}

type DB struct {
	db  *gorm.DB
	err error
}

func (d *DB) GetSTMPModel() *STMPModel {
	return NewSTMPModel(d.db, "t_STMP")
}
func (d *DB) GetMailModel() *MailModel {
	return NewMailModel(d.db, "t_mail")
}
func (d *DB) GetTokenModel() *TokenModel {
	return NewTokenModel(d.db, "t_token")
}
func CreateDB() *DB {
	return &DB{}
}
func (d *DB) Init(config *util.Config) error {
	var err error
	dbType := config.GetString("core", "db-type")
	if util.EqualsAnyIgnoreCase(dbType, "sqlite") {
		dbName := config.GetStringOrDefault("sqlite", "filename", "d-mail.db")
		d.db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
		if err != nil {
			d.err = err
			return err
		}
		d.err = err
		return err
	} else if util.EqualsAnyIgnoreCase(dbType, "mysql") {
		username := config.GetString("core", "username")
		password := config.GetString("core", "password")
		host := config.GetString("core", "host")
		port := config.GetString("core", "port")
		dbname := config.GetString("core", "dbname")
		charset := config.GetStringOrDefault("core", "charset", "UTF-8")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", username, password, host, port, dbname, charset)
		d.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
		if err != nil {
			d.err = err
			return err
		}
		d.err = err
		return err
	}
	d.err = NoDatabase
	return d.err
}

type NoDatabaseError struct {
	error
}

func (error *NoDatabaseError) Error() string {
	return "No database type selected"
}
