package main

import (
	"github.com/chuccp/d-mail/api"
	"github.com/chuccp/d-mail/core"
	"github.com/chuccp/d-mail/manage"
)

func main() {

	dMail := core.Create()
	dMail.AddServer(manage.NewServer())
	dMail.AddServer(api.NewServer())
	dMail.Start()

}
