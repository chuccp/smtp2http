package manage

import (
	"github.com/chuccp/smtp2http/core"
)

type User struct {
	context *core.Context
}

func (u *User) Init(context *core.Context, server core.IHttpServer) {
	server.SignIn("/signIn")
	server.Logout("/logout")
}
