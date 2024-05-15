package api

import "github.com/chuccp/d-mail/core"

type api interface {
	Init(context *core.Context)
}

type Server struct {
	apis    []api
	context *core.Context
}

func (s *Server) Name() string {
	return "api"
}
func (s *Server) Start() {
	for _, a := range s.apis {
		a.Init(s.context)
	}
}
func (s *Server) addApi(api api) {
	s.apis = append(s.apis, api)
}
func (s *Server) Init(context *core.Context) {
	s.context = context
	s.apis = make([]api, 0)
	s.addApi(&Stmp{})
	s.addApi(&Set{})
}
