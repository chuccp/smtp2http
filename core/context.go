package core

import (
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"runtime/debug"
)

type Context struct {
	db     *db.DB
	config *util.Config
	engine *gin.Engine
	log    *zap.Logger
}

func (c *Context) POST(relativePath string, handlers ...gin.HandlerFunc) {
	c.engine.POST(relativePath, handlers...)
}
func (c *Context) GetLog() *zap.Logger {
	return c.log
}
func (c *Context) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	c.engine.PUT(relativePath, handlers...)
}
func (c *Context) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	c.engine.DELETE(relativePath, handlers...)
}
func (c *Context) GET(relativePath string, handlers ...gin.HandlerFunc) {
	c.engine.GET(relativePath, handlers...)
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
