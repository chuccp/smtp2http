package schedule

import (
	"github.com/chuccp/smtp2http/core"
	"github.com/chuccp/smtp2http/db"
	"github.com/chuccp/smtp2http/service"
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
	lock       *sync.Mutex
	request    *util.Request
	log        *service.Log
	token      *service.Token
	context    *core.Context
}

func NewServer() *Server {
	return &Server{
		lock:       new(sync.Mutex),
		request:    util.NewRequest(),
		cronManage: newCronManage(),
	}
}

func (server *Server) Init(context *core.Context) {
	server.context = context
	server.log = service.NewLog(context.GetDb())
	server.token = service.NewToken(context.GetDb())
	context.SetSchedule(server)
}
func (server *Server) Name() string {
	return "schedule"
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
		err := server.Run(schedule)
		if err != nil {
			server.context.GetLog().Error("SendAPIMail Run log error", zap.Error(err))
		}
	}
}
func (server *Server) Run(schedule *db.Schedule) error {
	entryID, err := server.cronManage.cron.AddFunc(schedule.Cron, func() {
		byToken, err := server.token.GetOneByToken(schedule.Token)
		if err == nil {
			body, err := smtp.SendAPIMail(schedule, byToken.SMTP, byToken.ReceiveEmails)
			if err != nil {
				err := server.log.ContentError(byToken.SMTP, byToken.ReceiveEmails, schedule.Token, schedule.Name, body, err)
				if err != nil {

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
	isInit := server.context.GetConfig().GetBooleanOrDefault("core", "init", false)
	if isInit {
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
}
