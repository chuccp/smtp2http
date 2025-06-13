package schedule

import "github.com/robfig/cron/v3"

type CronManage struct {
	cronMap map[uint]cron.EntryID
	cron    *cron.Cron
}

func NewCronManage() *CronManage {
	return &CronManage{
		cronMap: make(map[uint]cron.EntryID),
		cron:    cron.New(cron.WithSeconds()),
	}
}

func (cronM *CronManage) Add(id uint, cron string, f func()) error {
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
