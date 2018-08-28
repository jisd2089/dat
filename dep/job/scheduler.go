package job

import (
	"time"

	"github.com/robfig/cron"
)

type SchedulerJob cron.Job

func Addfunc(express string, cmd func()) error {
	return scheduler.AddFunc(express, cmd)
}

func Schedule(duration time.Duration, cmd SchedulerJob) {
	scheduler.Schedule(cron.Every(duration), cmd)
}

func AddJob(spec string, cmd SchedulerJob) error {
	return scheduler.AddJob(spec, cmd)
}

func Stop() {
	scheduler.Stop()
}
