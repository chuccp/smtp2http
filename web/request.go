package web

import (
	"encoding/json"
	auth "github.com/abbot/go-http-auth"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type HandlerFunc func(req *Request) (any, error)

type Request struct {
	context     *gin.Context
	authRequest *auth.AuthenticatedRequest
}

func NewRequest(context *gin.Context, authRequest *auth.AuthenticatedRequest) *Request {
	return &Request{context: context, authRequest: authRequest}
}
func (r *Request) GetAuthRequest() *auth.AuthenticatedRequest {
	return r.authRequest
}

func (r *Request) ShouldBindBodyWithJSON(obj any) error {
	return r.context.ShouldBindBodyWith(obj, binding.JSON)
}

func (r *Request) GetContext() *gin.Context {
	return r.context
}

func (r *Request) GetAuthUsername() string {
	if r.authRequest != nil {
		return r.authRequest.Username
	}
	return ""
}

func (r *Request) ReadAuthUsername() string {
	authorization := r.context.Request.Header.Get("Authorization")
	if len(authorization) == 0 {
		return ""
	}
	data := auth.DigestAuthParams(authorization)
	if data == nil {
		return ""
	}
	return data["username"]
}
func (r *Request) ReadAuthUsernameAndCode() (string, string) {
	un := r.ReadAuthUsername()
	return GetUsernameAndCode(un)
}

func GetUsernameAndCode(user string) (string, string) {
	vs := strings.SplitN(user, "@", 2)
	if len(vs) > 1 {
		return vs[0], vs[1]
	}
	return user, ""
}

func (r *Request) FormValue(key string) string {
	return r.context.Request.FormValue(key)
}
func (r *Request) FormIntValue(key string) int {
	v := r.FormValue(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}
func (r *Request) FormInt64Value(key string) int64 {
	v := r.FormValue(key)
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}
	return i
}
func (r *Request) GetRemoteAddress() string {
	address := r.context.Request.RemoteAddr
	index := strings.Index(address, "_")
	if index > 0 {
		return address[:index]
	}
	return address
}
func (r *Request) GetPage() *Page {
	var page Page
	page.PageNo = r.FormIntValue("pageNo")
	page.PageSize = r.FormIntValue("pageSize")
	return &page
}
func (r *Request) GetRawRequest() *http.Request {
	return r.context.Request
}
func (r *Request) GetHeader(key string) string {
	return r.context.Request.Header.Get(key)
}
func (r *Request) GetResponseWriter() http.ResponseWriter {
	return r.context.Writer
}

func (r *Request) BodyJson(v any) ([]byte, error) {
	body, err := io.ReadAll(r.context.Request.Body)
	if err != nil {
		return body, err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return body, err
	}
	return body, nil
}

func (r *Request) FormFile(name string) (*multipart.FileHeader, error) {
	return r.context.FormFile(name)
}

func (r *Request) MultipartForm() (*multipart.Form, error) {
	return r.context.MultipartForm()
}
func (r *Request) IsForm() bool {
	return strings.Contains(r.context.GetHeader("Content-Type"), "application/x-www-form-urlencoded")
}
func (r *Request) IsMultipartForm() bool {
	return strings.Contains(r.context.GetHeader("Content-Type"), "multipart/form-data")
}

func (r *Request) Param(key string) string {
	return r.context.Param(key)
}

func (r *Request) BasicAuth() (username, password string, ok bool) {
	return r.context.Request.BasicAuth()
}
func (r *Request) Header(key, value string) {
	r.context.Header(key, value)
}
func (r *Request) Status(code int) {
	r.context.Status(code)
}
func (r *Request) String(code int, format string, values ...any) {
	r.context.String(code, format, values...)
}

func ToGinHandlerFuncs(handlers []HandlerFunc) []gin.HandlerFunc {
	var handlerFunc = make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handlerFunc[i] = ToGinHandlerFunc(handler)
	}
	return handlerFunc
}
func ToGinHandlerFunc(handler HandlerFunc) gin.HandlerFunc {
	handlerFunc := func(context *gin.Context) {
		value, err := handler(NewRequest(context, nil))
		if err != nil {
			context.Status(500)
			context.Writer.Write([]byte(err.Error()))
			context.Abort()
		} else {
			if value != nil {
				switch t := value.(type) {
				case string:
					context.Writer.Write([]byte(t))
				case *File:
					context.FileAttachment(t.GetPath(), t.GetFilename())
				default:
					context.AbortWithStatusJSON(200, t)
				}
			}
		}
	}
	return handlerFunc
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}
