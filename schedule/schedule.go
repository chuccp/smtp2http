package schedule

import (
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/smtp"
	"github.com/chuccp/smtp2http/util"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
)

type cronManage struct {
	cronMap map[uint]cron.EntryID
	cron    *cron.Cron
	lock    *sync.Mutex
}

func newCronManage() *cronManage {
	return &cronManage{
		cronMap: make(map[uint]cron.EntryID),
		cron:    cron.New(cron.WithSeconds()),
		lock:    new(sync.Mutex),
	}
}
func (cronM *cronManage) Start() {
	cronM.cron.Start()
}
func (cronM *cronManage) stop() {
	cronM.cron.Stop()
}

type Server struct {
	cronManage *cronManage
	context    *core.Context
	lock       *sync.Mutex
	request    *util.Request
}

func NewServer() *Server {
	return &Server{
		lock:    new(sync.Mutex),
		request: util.NewRequest(),
	}
}

func (server *Server) Init(context *core.Context) {
	server.context = context
}
func (server *Server) Name() string {
	return "cronManage"
}
func (server *Server) init() {
	if server.context.GetDb() == nil {
		return
	}
	schedules, err := server.context.GetDb().GetScheduleModel().FindAllByUse()
	if err != nil {
		return
	}
	for _, schedule := range schedules {
		entryID, err := server.cronManage.cron.AddFunc(schedule.Cron, func() {
			_, err := server.request.CallApi(schedule.Url, nil, schedule.Method, []byte(schedule.Body))
			if err != nil {
				return
			}
		})
		if err != nil {
			server.context.GetLog().Error("cron start error", zap.String("cron", schedule.Cron), zap.Error(err))
		} else {
			server.cronManage.cronMap[schedule.GetId()] = entryID
		}
	}
}
func (server *Server) Run(schedule *db.Schedule) error {
	entryID, err := server.cronManage.cron.AddFunc(schedule.Cron, func() {
		byToken, err := server.context.GetTokenService().GetOneByToken(schedule.Token)
		if err == nil {
			body, err := smtp.SendAPIMail(schedule, byToken.SMTP, byToken.ReceiveEmails)
			if err != nil {
				err := server.context.GetLogService().ContentError(byToken.SMTP, byToken.ReceiveEmails, schedule.Token, schedule.Name, body, err)
				if err != nil {
					server.context.GetLog().Error("SendAPIMail log error", zap.Error(err))
				}
			} else {
				err := server.context.GetLogService().ContentSuccess(byToken.SMTP, byToken.ReceiveEmails, schedule.Token, schedule.Name, body)
				if err != nil {
					server.context.GetLog().Error("SendAPIMail log error", zap.Error(err))
				}
			}
		}
	})
	if err != nil {
		server.context.GetLog().Error("cron start error", zap.String("cron", schedule.Cron), zap.Error(err))
	} else {
		server.cronManage.cronMap[schedule.GetId()] = entryID
	}
	return err
}

func (server *Server) Start() {
	server.lock.Lock()
	defer server.lock.Unlock()
	if server.cronManage != nil {
		server.cronManage.stop()
	}
	server.cronManage = newCronManage()
	server.init()
	server.cronManage.Start()
	server.context.GetLog().Info("start Schedule")
}
