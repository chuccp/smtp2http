package main

import (
	"github.com/chuccp/d-mail/api"
	"github.com/chuccp/d-mail/core"
)

func main() {

	dMail := core.Create()
	dMail.AddServer(&api.Server{})
	dMail.Start()

}
