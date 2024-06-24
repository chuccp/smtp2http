package core

import (
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/entity"
	"github.com/chuccp/d-mail/util"
	"github.com/chuccp/d-mail/web"
	"go.uber.org/zap"
	"log"
	"runtime/debug"
)

type Context struct {
	db         *db.DB
	config     *util.Config
	log        *zap.Logger
	httpServer *util.HttpServer
}

func (c *Context) GetLog() *zap.Logger {
	return c.log
}
func (c *Context) GetDb() *db.DB {
	return c.db
}
func (c *Context) GetConfig() *util.Config {
	return c.config
}
func (c *Context) IsInit() bool {
	return c.config.GetBooleanOrDefault("core", "init", false)
}

func (c *Context) post(relativePath string, handlers ...web.HandlerFunc) {
	c.httpServer.POST(relativePath, handlers...)
}
func (c *Context) put(relativePath string, handlers ...web.HandlerFunc) {
	c.httpServer.PUT(relativePath, handlers...)
}
func (c *Context) delete(relativePath string, handlers ...web.HandlerFunc) {
	c.httpServer.DELETE(relativePath, handlers...)
}

func (c *Context) get(relativePath string, handlers ...web.HandlerFunc) {
	c.httpServer.GET(relativePath, handlers...)
}

func (c *Context) GetDefaultSetInfo() *entity.SetInfo {
	var setInfo entity.SetInfo
	setInfo.HasInit = c.config.GetBooleanOrDefault("core", "init", false)
	setInfo.DbType = c.config.GetString("core", "db-type")
	var sqlite entity.Sqlite
	sqlite.Filename = c.config.GetString("sqlite", "filename")
	setInfo.Sqlite = &sqlite
	var mysql entity.Mysql
	mysql.Host = c.config.GetString("mysql", "host")
	mysql.Port = c.config.GetIntOrDefault("mysql", "port", 0)
	mysql.Dbname = c.config.GetString("mysql", "dbname")
	mysql.Username = c.config.GetString("mysql", "username")
	mysql.Password = c.config.GetString("mysql", "password")
	mysql.Charset = c.config.GetString("mysql", "charset")
	setInfo.Mysql = &mysql
	return &setInfo
}

func (c *Context) UpdateSetInfo(setInfo *entity.SetInfo) error {
	c.config.SetBoolean("core", "init", true)
	c.config.SetString("core", "db-type", setInfo.DbType)
	if setInfo.DbType == "sqlite" {
		c.config.SetString("sqlite", "filename", setInfo.Sqlite.Filename)
	} else {
		c.config.SetString("mysql", "host", setInfo.Mysql.Host)
		c.config.SetInt("mysql", "port", setInfo.Mysql.Port)
		c.config.SetString("mysql", "dbname", setInfo.Mysql.Dbname)
		c.config.SetString("mysql", "charset", setInfo.Mysql.Charset)
		c.config.SetString("mysql", "username", setInfo.Mysql.Username)
		c.config.SetString("mysql", "password", setInfo.Mysql.Password)
	}
	err := c.config.Save()
	if err != nil {
		return err
	}
	err = c.initDb()
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) initDb() error {
	if c.IsInit() {
		_db_ := db.CreateDB()
		err := _db_.Init(c.config)
		if err != nil {
			return err
		}
		c.db = _db_

		c.db.GetSTMPModel().CreateTable()
		c.db.GetMailModel().CreateTable()
		c.db.GetTokenModel().CreateTable()
		c.db.GetLogModel().CreateTable()
	}
	return nil
}

func (c *Context) Go(handle func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				s := string(debug.Stack())
				log.Println(err)
				log.Println(s)
				c.log.Error("Go", zap.Any("err", err), zap.String("info", s))
			}
		}()
		handle()
	}()
}
func (c *Context) GetCfgInt(section string, name string) int {
	getInt, err := c.config.GetInt(section, name)
	if err != nil {
		return 0
	} else {
		return getInt
	}
}

func (c *Context) GetCfgString(section string, name string) string {
	return c.config.GetString(section, name)
}
