package schedule

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type CronManage struct {
	cronMap map[uint]cron.EntryID
	cron    *cron.Cron
	lock    *sync.Mutex
}

func NewCronManage() *CronManage {
	return &CronManage{
		cronMap: make(map[uint]cron.EntryID),
		cron:    cron.New(cron.WithSeconds()),
		lock:    new(sync.Mutex),
	}
}

func (cronM *CronManage) Add(id uint, cron string, f func()) error {
	cronM.lock.Lock()
	defer cronM.lock.Unlock()
	entryId, err := cronM.cron.AddFunc(cron, f)
	if err != nil {
		return err
	}
	cronM.cronMap[id] = entryId
	return nil
}
func (cronM *CronManage) Start() {
	cronM.cron.Start()
}
