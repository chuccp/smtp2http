package main

import (
	"github.com/chuccp/smtp2http/api"
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/manage"
)

func main() {

	dMail := core.Create()
	dMail.AddServer(manage.NewServer())
	dMail.AddServer(api.NewServer())
	dMail.Start()

}
