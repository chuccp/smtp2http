package api

import (
	"github.com/chuccp/d-mail/core"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Stmp struct {
	context *core.Context
}

func (stmp *Stmp) getOne(context *gin.Context) {

	id := context.Param("id")
	stmp.context.GetLog().Info("getOne", zap.String("id", id))

}
func (stmp *Stmp) getList(context *gin.Context) {

}
func (stmp *Stmp) Init(context *core.Context) {
	stmp.context = context
	context.GET("/stmp/:id", stmp.getOne)
	context.GET("/stmp", stmp.getList)
}
