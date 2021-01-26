//go:generate pkger
package main

import (
	"fmt"
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/Hive-Gay/supreme-robot/webapp"
	"github.com/garyburd/redigo/redis"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/thatisuday/commando"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var logger *loggo.Logger

func main() {

	// Init Logging
	newLogger := loggo.GetLogger("main")
	logger = &newLogger

	_, err := loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	if err != nil {
		logger.Errorf("Error configuring Color Logger: %s", err.Error())
		return
	}

	// configure commando
	commando.
		SetExecutableName("supreme-robot").
		SetVersion("0.0.1").
		SetDescription("This tool runs hive operations.")

	commando.
		Register(nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			commando.Parse([]string{"help"})
		})

	commando.
		Register("server").
		SetShortDescription("runs the supreme robot server").
		AddFlag("log,l", "level of logging", commando.String, "info").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			logLevel, _ := flags["log"].GetString()

			err := loggo.ConfigureLoggers(fmt.Sprintf("<root>=%s", strings.ToUpper(logLevel)))
			if err != nil {
				logger.Errorf("could not configure logger: %s", err.Error())
				return
			}

			logger.Infof("Starting Supreme Robot")
			redisPool := initRedisPool()

		err = models.Init()
		if err != nil {
			logger.Errorf("could not start models: %s", err.Error())
			return
		}

			err = webapp.Init(redisPool)
			if err != nil {
				logger.Errorf("could not start webapp: %s", err.Error())
				return
			}


			// Wait for SIGINT and SIGTERM (HIT CTRL-C)
			nch := make(chan os.Signal)
			signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
			logger.Infof("%s", <-nch)

		})

	commando.Parse(nil)

}

func initRedisPool() *redis.Pool {
	logger.Debugf("starting redis pool")

	return &redis.Pool{
		MaxActive: 10,
		MaxIdle:   10,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}
}
