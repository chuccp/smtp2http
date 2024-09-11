package main

import (
	"flag"
	"github.com/chuccp/smtp2http/api"
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/manage"
)

func main() {

	var webPort int
	var apiPort int
	flag.IntVar(&webPort, "web_port", 0, "web port")
	flag.IntVar(&apiPort, "api_port", 0, "api port")
	flag.Parse()
	dMail := core.Create()
	dMail.AddServer(manage.NewServer())
	dMail.AddServer(api.NewServer())
	dMail.Start(webPort, apiPort)

}
