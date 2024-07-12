package manage

import (
	"github.com/chuccp/d-mail/core"
)

type User struct {
	context *core.Context
}

func (u *User) Init(context *core.Context, server core.IHttpServer) {
	server.SignIn("/signIn")
}
