package core

type Server interface {
	Start()
	Init(context *Context)
	Name() string
}
