package api

import (
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/entity"
	"github.com/chuccp/d-mail/web"
	"go.uber.org/zap"
)

type Set struct {
	context *core.Context
}

func (set *Set) putSet(req *web.Request) (any, error) {
	var setInfo entity.SetInfo
	err := req.ShouldBindBodyWithJSON(&setInfo)
	if err != nil {
		return nil, err
	} else {
		set.context.GetLog().Debug("putSet", zap.Any("setInfo", &setInfo))
		err := set.context.UpdateSetInfo(&setInfo)
		if err != nil {
			return nil, err
		}
		return "ok", nil
	}
}
func (set *Set) getSet(req *web.Request) (any, error) {
	return &entity.System{HasInit: set.context.IsInit(), HasLogin: false}, nil
}
func (set *Set) defaultSet(req *web.Request) (any, error) {
	return set.context.GetDefaultSetInfo(), nil
}

func (set *Set) Init(context *core.Context) {
	set.context = context
	context.GET("/set", set.getSet)
	context.GET("/defaultSet", set.defaultSet)
	context.PUT("/set", set.putSet)
}
