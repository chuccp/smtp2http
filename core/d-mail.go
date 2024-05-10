package core

import (
	"errors"
	"github.com/chuccp/d-mail/db"
	"github.com/chuccp/d-mail/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

type DMail struct {
	context *Context
	engine  *gin.Engine
	log     *zap.Logger
	config  *util.Config
	db      *db.DB
	servers []Server
}

func Create() *DMail {
	return &DMail{servers: make([]Server, 0)}
}
func (m *DMail) AddServer(server Server) {
	m.servers = append(m.servers, server)
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
	m.db = db.CreateDB()
	err = m.db.Init(m.config)
	if err != nil {
		ok := errors.Is(err, db.NoDatabase)
		if !ok {
			log.Panic(err)
			return
		}
	}
	m.engine = gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // 允许的域名列表，可以使用 * 来允许所有域名
	config.AllowHeaders = []string{"*"} // 允
	m.engine.Use(cors.New(config))
	port, err := configure.GetInt("core", "port")
	if err != nil {
		m.log.Panic("Start", zap.Error(err))
		return
	}
	m.context = &Context{log: m.log, engine: m.engine, config: m.config}

	for _, server := range m.servers {
		server.Init(m.context)
		m.context.Go(server.Start)
	}

	err = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), m.engine)
	if err != nil {
		m.log.Panic("Start", zap.Error(err))
		return
	}
}
