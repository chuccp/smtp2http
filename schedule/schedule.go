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
	"time"
)

type cronManage struct {
	cronMap map[uint]cron.EntryID
	cron    *cron.Cron
	lock    *sync.Mutex
	context *core.Context
	log     *service.Log
	token   *service.Token
	isStart bool
	isStop  bool
}

func newCronManage(context *core.Context) *cronManage {
	return &cronManage{
		cronMap: make(map[uint]cron.EntryID),
		cron:    cron.New(cron.WithSeconds()),
		lock:    new(sync.Mutex),
		context: context,
		isStart: false,
		isStop:  false,
		log:     service.NewLog(context.GetDb()),
		token:   service.NewToken(context.GetDb()),
	}
}
func (cronM *cronManage) Start() {
	if cronM.context.GetDb() == nil {
		return
	}
	cronM.isStart = true
	cronM.cron.Start()
	go cronM.run()
}
func (cronM *cronManage) run() {
	for {
		if cronM.isStop {
			break
		}
		schedules, err := cronM.context.GetDb().GetScheduleModel().FindAllByUse()
		if err != nil {
			return
		}
		runIds := make([]uint, 0)
		for _, schedule := range schedules {
			_, ok := cronM.cronMap[schedule.GetId()]
			if ok {
				continue
			}
			runIds = append(runIds, schedule.GetId())
			cronM.addSchedule(schedule)
		}
		removeIds := make([]uint, 0)
		for key, value := range cronM.cronMap {
			if !util.ContainsNumberAny(key, runIds...) {
				cronM.cron.Remove(value)
				removeIds = append(removeIds, key)
			}
		}
		for _, key := range removeIds {
			delete(cronM.cronMap, key)
		}
		time.Sleep(time.Second * 10)
	}
}

func (cronM *cronManage) addSchedule(schedule *db.Schedule) {

	entryID, err := cronM.cron.AddFunc(schedule.Cron, func() {
		byToken, err := cronM.token.GetOneByToken(schedule.Token)
		if err == nil {
			body, err := smtp.SendAPIMail(schedule, byToken.SMTP, byToken.ReceiveEmails)
			if err != nil {
				err := cronM.log.ContentError(byToken.SMTP, byToken.ReceiveEmails, schedule.Token, schedule.Name, body, err)
				if err != nil {

				}
			} else {
				err := cronM.context.GetLogService().ContentSuccess(byToken.SMTP, byToken.ReceiveEmails, schedule.Token, schedule.Name, body)
				if err != nil {
					cronM.context.GetLog().Error("SendAPIMail log error", zap.Error(err))
				}
			}
		}
	})
	if err != nil {
		cronM.context.GetLog().Error("cron start error", zap.String("cron", schedule.Cron), zap.Error(err))
	} else {
		cronM.cronMap[schedule.GetId()] = entryID
	}
}

func (cronM *cronManage) stop() {
	if !cronM.isStart {
		return
	}
	cronM.isStop = true
	cronM.cron.Stop()
}

type Server struct {
	cronManage *cronManage
	lock       *sync.Mutex
	request    *util.Request
	context    *core.Context
}

func NewServer() *Server {
	return &Server{
		lock:    new(sync.Mutex),
		request: util.NewRequest(),
	}
}

func (server *Server) Init(context *core.Context) {
	context.SetSchedule(server)
	server.context = context
}
func (server *Server) Name() string {
	return "schedule"
}
func (server *Server) Run(schedule *db.Schedule) error {
	return nil
}
func (server *Server) Stop(Id uint) error {
	return nil
}

func (server *Server) Start() {
	isInit := server.context.GetConfig().GetBooleanOrDefault("core", "init", false)
	if isInit {
		server.lock.Lock()
		defer server.lock.Unlock()
		if server.cronManage != nil {
			server.cronManage.stop()
		}
		server.cronManage = newCronManage(server.context)
		server.cronManage.Start()
		server.context.GetLog().Info("start Schedule")
	}
}
