package util

import (
	"github.com/chuccp/d-mail/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func NewServer() *HttpServer {
	engine := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // 允许的域名列表，可以使用 * 来允许所有域名
	config.AllowHeaders = []string{"*"} // 允
	engine.Use(cors.New(config))
	return &HttpServer{engine: engine}
}

type HttpServer struct {
	isTls  bool
	engine *gin.Engine
}

func (hs *HttpServer) IsTls() bool {
	return hs.isTls
}
func (hs *HttpServer) DELETE(pattern string, handlers ...web.HandlerFunc) {
	hs.engine.DELETE(pattern, web.ToGinHandlerFuncs(handlers)...)
}
func (hs *HttpServer) PUT(pattern string, handlers ...web.HandlerFunc) {
	hs.engine.PUT(pattern, web.ToGinHandlerFuncs(handlers)...)
}

func (hs *HttpServer) POST(pattern string, handlers ...web.HandlerFunc) {
	hs.engine.POST(pattern, web.ToGinHandlerFuncs(handlers)...)
}
func (hs *HttpServer) GET(pattern string, handlers ...web.HandlerFunc) {
	hs.engine.GET(pattern, web.ToGinHandlerFuncs(handlers)...)
}

const MaxHeaderBytes = 8192

const MaxReadHeaderTimeout = time.Second * 30

const MaxReadTimeout = time.Minute * 10

func (hs *HttpServer) Start(port int) error {
	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           hs.engine,
		ReadHeaderTimeout: MaxReadHeaderTimeout,
		MaxHeaderBytes:    MaxHeaderBytes,
		ReadTimeout:       MaxReadTimeout,
	}
	hs.isTls = false
	error := srv.ListenAndServe()
	return error
}
func (hs *HttpServer) StartTLS(port int, certFile, keyFile string) error {
	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           hs.engine,
		ReadHeaderTimeout: MaxReadHeaderTimeout,
		MaxHeaderBytes:    MaxHeaderBytes,
		ReadTimeout:       MaxReadTimeout,
	}
	hs.isTls = true
	return srv.ListenAndServeTLS(certFile, keyFile)
}
func (hs *HttpServer) StartAutoTLS(port int, certFile, keyFile string) error {
	if len(certFile) > 0 && len(keyFile) > 0 {
		return hs.StartTLS(port, certFile, keyFile)
	} else {
		return hs.Start(port)
	}
}
