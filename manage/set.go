package manage

import (
	"github.com/chuccp/d-mail/config"
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/web"
	"go.uber.org/zap"
)

type Set struct {
	context *core.Context
}

func (set *Set) putSet(req *web.Request) (any, error) {
	var setInfo config.SetInfo
	err := req.ShouldBindBodyWithJSON(&setInfo)
	if err != nil {
		return nil, err
	} else {
		set.context.GetLog().Debug("putSet", zap.Any("setInfo", &setInfo))
		setInfo.HasInit = true
		err := set.context.UpdateSetInfo(&setInfo)
		if err != nil {
			return nil, err
		}
		return "ok", nil
	}
}
func (set *Set) getSet(req *web.Request) (any, error) {
	hasLogin := req.GetDigestAuth().HasSign(req.GetContext())
	return &config.System{HasInit: set.context.IsInit(), HasLogin: hasLogin}, nil
}
func (set *Set) defaultSet(req *web.Request) (any, error) {
	return set.context.GetDefaultSetInfo(), nil
}

func (set *Set) Init(context *core.Context, server core.IHttpServer) {
	set.context = context
	server.GET("/set", set.getSet)
	server.GET("/defaultSet", set.defaultSet)
	server.PUT("/set", set.putSet)
}
