package core

import (
	"github.com/chuccp/d-mail/util"
	"go.uber.org/zap"
	"log"
)

type DMail struct {
	context    *Context
	httpServer *util.HttpServer
	log        *zap.Logger
	config     *util.Config
	servers    []Server
}

func Create() *DMail {
	return &DMail{servers: make([]Server, 0), httpServer: util.NewServer()}
}
func (m *DMail) AddServer(server Server) {
	m.servers = append(m.servers, server)
}

func (m *DMail) startHttpServer() error {
	port := m.context.GetCfgInt("manage", "port")
	certFile := m.context.GetCfgString("manage", "certFile")
	keyFile := m.context.GetCfgString("manage", "keyFile")
	m.context.log.Info("startHttpServer", zap.String("name", "manage"), zap.Int("port", port))
	err := m.httpServer.StartAutoTLS(port, certFile, keyFile)
	if err != nil {
		m.context.log.Error("服务启动失败", zap.String("name", "DMail"), zap.Int("port", port), zap.Error(err))
		return err
	}
	return nil
}

func (m *DMail) Start() {
	configure, err := util.LoadFile("config.ini")
	if err != nil {
		log.Panic(err)
		return
	}
	m.config = configure
	logPath := configure.GetStringOrDefault("log", "filename", "run.log")
	m.log, err = initLogger(logPath)
	if err != nil {
		log.Panic(err)
		return
	}
	m.context = &Context{log: m.log, httpServer: m.httpServer, config: m.config}
	isInit := configure.GetBooleanOrDefault("core", "init", false)
	if isInit {
		err := m.context.initDb()
		if err != nil {
			m.log.Panic("initDb", zap.Error(err))
			return
		}
	}
	for _, server := range m.servers {
		if s, ok := server.(IHttpServer); ok {
			s.init(m.context)
		}
		server.Init(m.context)
		if s, ok := server.(IHttpServer); ok {
			if !s.useCorePort() {
				go func() {
					err := s.start()
					if err != nil {
						log.Panic(err)
					}
				}()
			}
		}
	}
	err = m.startHttpServer()
	if err != nil {
		m.log.Panic("Start", zap.Error(err))
		return
	}
}
