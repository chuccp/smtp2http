package web

import (
	"github.com/chuccp/d-mail/login"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func NewServer(digestAuth *login.DigestAuth) *HttpServer {
	engine := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // 允许的域名列表，可以使用 * 来允许所有域名
	config.AllowHeaders = []string{"*"} // 允
	engine.Use(cors.New(config))
	return &HttpServer{engine: engine, digestAuth: digestAuth}
}

type HttpServer struct {
	isTls      bool
	engine     *gin.Engine
	digestAuth *login.DigestAuth
}

func (hs *HttpServer) IsTls() bool {
	return hs.isTls
}
func (hs *HttpServer) DELETE(pattern string, handlers ...HandlerFunc) {
	hs.engine.DELETE(pattern, ToGinHandlerFuncs(handlers, hs.digestAuth)...)
}
func (hs *HttpServer) PUT(pattern string, handlers ...HandlerFunc) {
	hs.engine.PUT(pattern, ToGinHandlerFuncs(handlers, hs.digestAuth)...)
}
func (hs *HttpServer) Any(pattern string, handlers ...HandlerFunc) {
	hs.engine.Any(pattern, ToGinHandlerFuncs(handlers, hs.digestAuth)...)
}

func (hs *HttpServer) POST(pattern string, handlers ...HandlerFunc) {
	hs.engine.POST(pattern, ToGinHandlerFuncs(handlers, hs.digestAuth)...)
}
func (hs *HttpServer) GET(pattern string, handlers ...HandlerFunc) {
	hs.engine.GET(pattern, ToGinHandlerFuncs(handlers, hs.digestAuth)...)
}
func (hs *HttpServer) SignIn(pattern string, handlers ...HandlerFunc) {
	hs.engine.GET(pattern, ToGinHandlerFuncs(handlers, hs.digestAuth)...)
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
