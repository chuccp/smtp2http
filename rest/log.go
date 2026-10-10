package rest

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	wf "github.com/chuccp/go-web-frame"
	auth2 "github.com/chuccp/go-web-frame/component/auth"
	"github.com/chuccp/go-web-frame/core"
	"github.com/chuccp/go-web-frame/log"
	"github.com/chuccp/go-web-frame/web"
	"github.com/chuccp/http2smtp/auth"
	"github.com/chuccp/http2smtp/entity"
	"github.com/chuccp/http2smtp/model"
	"github.com/chuccp/http2smtp/service"
	"github.com/chuccp/http2smtp/smtp"
	"go.uber.org/zap"
)

type Log struct {
	context     *core.Context
	logModel    *model.LogModel
	userService *service.UserService
}

// The log page is the instance's whole send record for an admin and the caller's own
// entries for everyone else, matching the other lists: every query below adds the owner
// filter only for a non-admin, and only then is a miss reported as "not found" rather
// than "forbidden" so existence does not leak.

func (l *Log) getOne(req *web.Request) (any, error) {
	id := req.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	user, err := auth.User(req, l.context)
	if user == nil {
		return nil, err
	}
	query := l.logModel.Query().Where("id = ?", uint(atoi))
	if !user.IsAdmin {
		query = query.Where("user_id = ?", user.Id)
	}
	one, err := query.One()
	if err != nil {
		return nil, err
	}
	if one != nil {
		one.StatusStr = entity.StatusText(one.Status)
	}
	return one, nil
}

func (l *Log) getPage(req *web.Request) (any, error) {
	page, err := req.Page()
	if err != nil {
		return nil, err
	}
	user, err := auth.User(req, l.context)
	if user == nil {
		return nil, err
	}
	searchKey := req.GetFormParam("searchKey")
	adminOnly := req.Query("adminOnly") == "true"
	query := l.logModel.Query()
	if user.IsAdmin {
		if adminOnly {
			query = query.Where("user_id IN (SELECT id FROM t_user WHERE is_admin = ?)", true)
		}
	} else {
		query = query.Where("user_id = ?", user.Id)
	}
	if searchKey != "" {
		query = query.Where("(name LIKE ? OR mail LIKE ? OR subject LIKE ?)", "%"+searchKey+"%", "%"+searchKey+"%", "%"+searchKey+"%")
	}
	result, err := query.Order("id desc").PageForWeb(page)
	if err != nil {
		return nil, err
	}
	for _, one := range result.List {
		if one != nil {
			one.StatusStr = entity.StatusText(one.Status)
		}
	}
	if user.IsAdmin {
		userIds := make([]uint, 0)
		for _, one := range result.List {
			if one != nil && one.UserId > 0 {
				userIds = append(userIds, one.UserId)
			}
		}
		l.userService.FillUserNames(userIds, func(uid uint, name string) {
			for _, one := range result.List {
				if one != nil && one.UserId == uid {
					one.UserName = name
				}
			}
		})
	}
	return result, nil
}

// downLoad serves an attachment recorded on a log the caller may read: their own, or any
// entry when they are an admin. The path is resolved from the log record, never from
// client input.
func (l *Log) downLoad(req *web.Request) (any, error) {
	logId, err := strconv.Atoi(req.GetFormParam("logId"))
	if err != nil {
		return nil, errors.New("logId is required")
	}
	fileName := req.GetFormParam("file")
	if fileName == "" {
		return nil, errors.New("file is required")
	}
	user, err := auth.User(req, l.context)
	if user == nil {
		return nil, err
	}
	query := l.logModel.Query().Where("id = ?", uint(logId))
	if !user.IsAdmin {
		query = query.Where("user_id = ?", user.Id)
	}
	one, err := query.One()
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, errors.New("log not found")
	}
	filePath, err := FilePathOf(one.Files, fileName)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(filePath); err != nil {
		return nil, errors.New("attachment file is no longer available")
	}
	log.Info("downLoad", zap.Uint("logId", uint(logId)), zap.String("file", fileName), zap.Uint("userId", user.Id))
	return web.CreateFileResponse(filePath), nil
}

// FilePathOf resolves the stored path of a named attachment from a log's files JSON.
func FilePathOf(filesJSON string, name string) (string, error) {
	if filesJSON == "" {
		return "", errors.New("log has no attachments")
	}
	var files []*smtp.File
	if err := json.Unmarshal([]byte(filesJSON), &files); err != nil {
		return "", errors.New("invalid attachment record")
	}
	for _, file := range files {
		if file != nil && file.Name == name && file.FilePath != "" {
			return file.FilePath, nil
		}
	}
	return "", errors.New("attachment not found")
}

func (l *Log) Init(context *core.Context) error {
	l.context = context
	l.logModel = wf.GetModel[*model.LogModel](context)
	l.userService = wf.GetService[*service.UserService](context)
	context.Get("/log/:id", l.getOne).WithMeta(auth2.WithLogin())
	context.Get("/log", l.getPage).WithMeta(auth2.WithLogin())
	context.Get("/download", l.downLoad).WithMeta(auth2.WithLogin())
	return nil

}
