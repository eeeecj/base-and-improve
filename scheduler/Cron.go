package scheduler

import "github.com/robfig/cron/v3"

func init() {
	Cron.Start()
}

var Cron = cron.New()
