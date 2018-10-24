package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/leangeder/gravitywell/configuration"
	"github.com/leangeder/gravitywell/scheduler"
	"os"
)

func ApplicationApply(config *configuration.ApplicationConfig, opt *configuration.Options) {
	log.Println("ApplicationApply called")

	sh, err := scheduler.NewScheduler(config)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	if err := sh.Run(*opt); err != nil {
		log.Warn(err.Error())
		os.Exit(1)
	}
}
