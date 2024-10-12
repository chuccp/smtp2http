package core

import (
	"errors"
	"github.com/chuccp/smtp2http/config"
	"github.com/chuccp/smtp2http/login"
	"github.com/chuccp/smtp2http/web"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type DMail struct {
	context     *Context
	httpServer  *web.HttpServer
	log         *zap.Logger
	config      *config.Config
	servers     []Server
	webPort     int
	apiPort     int
	iHttpServer []IHttpServer
}

func Create() *DMail {
	return &DMail{webPort: 0, apiPort: 0, servers: make([]Server, 0), config: config.NewConfig(), iHttpServer: make([]IHttpServer, 0)}
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
func (m *DMail) Start(webPort int, apiPort int) {
	m.webPort = webPort
	m.apiPort = apiPort
	for {
		m.reStart()
	}
}
func (m *DMail) ReStart() {
	for _, server := range m.iHttpServer {
		server.Stop()
	}
	time.Sleep(2 * time.Second)
	m.httpServer.Stop()
}
func (m *DMail) reStart() {
	m.iHttpServer = make([]IHttpServer, 0)
	err := m.config.Init(m.webPort, m.apiPort)
	if err != nil {
		log.Panic(err)
		return
	}
	logPath := m.config.GetStringOrDefault("log", "filename", "run.log")
	logger, err := initLogger(logPath)
	if err != nil {
		log.Panic(err)
		return
	}
	m.log = logger
	m.context = &Context{log: m.log, config: m.config, reStart: m.ReStart}
	digestAuth := login.NewDigestAuth(m.context.SecretProvider)
	m.context.digestAuth = digestAuth
	m.httpServer = web.NewServer(digestAuth)
	m.context.httpServer = m.httpServer
	isInit := m.config.GetBooleanOrDefault("core", "init", false)
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
					m.iHttpServer = append(m.iHttpServer, s)
					err := s.start()
					if !errors.Is(err, http.ErrServerClosed) && err != nil {
						log.Panic(err)
					}
				}()
			}
		}
	}
	err = m.startHttpServer()
	if !errors.Is(err, http.ErrServerClosed) && err != nil {
		m.log.Panic("Start", zap.Error(err))
		return
	}
}
