package core

import (
	"github.com/chuccp/d-mail/util"
	"github.com/chuccp/d-mail/web"
	"go.uber.org/zap"
)

type Server interface {
	Init(context *Context)
	Name() string
}
type IHttpServer interface {
	init(context *Context)
	start() error
	useCorePort() bool
	GET(relativePath string, handlers ...web.HandlerFunc)
	POST(relativePath string, handlers ...web.HandlerFunc)
	DELETE(relativePath string, handlers ...web.HandlerFunc)
	PUT(relativePath string, handlers ...web.HandlerFunc)
}
type httpServer struct {
	context    *Context
	port       int
	usePort    int
	httpServer *util.HttpServer
	certFile   string
	keyFile    string
	name       string
}

func NewHttpServer(name string) IHttpServer {
	return &httpServer{name: name}
}
func (server *httpServer) PUT(pattern string, handlers ...web.HandlerFunc) {
	if server.port > 0 {
		server.httpServer.PUT(pattern, handlers...)
	} else {
		server.context.put(pattern, handlers...)
	}
}

func (server *httpServer) POST(pattern string, handlers ...web.HandlerFunc) {
	if server.port > 0 {
		server.httpServer.POST(pattern, handlers...)
	} else {
		server.context.post(pattern, handlers...)
	}
}

func (server *httpServer) DELETE(pattern string, handlers ...web.HandlerFunc) {
	if server.port > 0 {
		server.httpServer.DELETE(pattern, handlers...)
	} else {
		server.context.delete(pattern, handlers...)
	}
}
func (server *httpServer) GET(pattern string, handlers ...web.HandlerFunc) {
	if server.port > 0 {
		server.httpServer.GET(pattern, handlers...)
	} else {
		server.context.get(pattern, handlers...)
	}
}

func (server *httpServer) IsTls() bool {
	if server.port > 0 {
		return server.httpServer.IsTls()
	} else {
		return server.context.httpServer.IsTls()
	}
}

func (server *httpServer) init(context *Context) {
	server.context = context
	port := context.GetCfgInt(server.name, "port")
	corePort := context.GetCfgInt("manage", "port")
	if port > 0 && corePort != port {
		context.log.Info("服务名称与端口", zap.String("name", server.name), zap.Int("port", port))
		server.certFile = context.GetCfgString(server.name, "certFile")
		server.keyFile = context.GetCfgString(server.name, "keyFile")
		server.port = port
		server.usePort = port
		server.httpServer = util.NewServer()
	} else {
		server.usePort = corePort
		context.log.Info("服务名称与端口", zap.String("name", server.name), zap.Int("port", corePort))
	}

}
func (server *httpServer) useCorePort() bool {
	return server.port < 1
}

func (server *httpServer) start() error {
	if server.port > 0 {
		err := server.httpServer.StartAutoTLS(server.port, server.certFile, server.keyFile)
		if err != nil {
			server.context.log.Error("服务启动失败", zap.String("name", server.name), zap.Int("port", server.port), zap.Error(err))
			return err
		}
		return nil
	} else {
		return nil
	}
	return nil
}
