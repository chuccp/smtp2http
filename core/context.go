package core

import (
	"github.com/chuccp/d-mail/config"
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/login"
	"github.com/chuccp/d-mail/util"
	"github.com/chuccp/d-mail/web"
	"go.uber.org/zap"
	"log"
	"runtime/debug"
)

type Context struct {
	db         *db.DB
	config     *config.Config
	log        *zap.Logger
	httpServer *web.HttpServer
	digestAuth *login.DigestAuth
}

func (c *Context) GetDigestAuth() *login.DigestAuth {
	return c.digestAuth
}
func (c *Context) GetLog() *zap.Logger {
	return c.log
}
func (c *Context) GetDb() *db.DB {
	return c.db
}

func (c *Context) SecretProvider(user string) string {
	password := c.config.GetString("manage", "password")
	username := c.config.GetString("manage", "username")
	if username == user {
		return util.MD5Str(util.MD5Str(password) + username)
	}
	return ""
}

func (c *Context) GetConfig() *config.Config {
	return c.config
}

func (c *Context) IsInit() bool {
	return c.config.GetBooleanOrDefault("core", "init", false)
}

func (c *Context) post(relativePath string, handlers ...web.HandlerFunc) {
	c.httpServer.POST(relativePath, handlers...)
}
func (c *Context) any(relativePath string, handlers ...web.HandlerFunc) {
	c.httpServer.Any(relativePath, handlers...)
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

func (c *Context) GetDefaultSetInfo() *config.SetInfo {
	return c.config.ReadSetInfo()
}

func (c *Context) UpdateSetInfo(setInfo *config.SetInfo) error {
	err := c.initDbBySetInfo(setInfo)
	if err != nil {
		return err
	}
	setInfo.HasInit = true
	err = c.config.UpdateSetInfo(setInfo)
	if err != nil {
		return err
	}
	return nil
}
func (c *Context) initDbBySetInfo(setInfo *config.SetInfo) error {
	_db_ := db.CreateDB()
	err := _db_.InitBySetInfo(setInfo)
	if err != nil {
		return err
	}
	c.db = _db_
	c.db.GetSTMPModel().CreateTable()
	c.db.GetMailModel().CreateTable()
	c.db.GetTokenModel().CreateTable()
	c.db.GetLogModel().CreateTable()
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
	return c.config.GetInt(section, name)
}

func (c *Context) GetCfgString(section string, name string) string {
	return c.config.GetString(section, name)
}
