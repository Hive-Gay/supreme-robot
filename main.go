//go:generate pkger
package main

import (
	"errors"
	"fmt"
	"github.com/Hive-Gay/supreme-robot/jobs"
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/Hive-Gay/supreme-robot/twilio"
	"github.com/Hive-Gay/supreme-robot/webapp"
	"github.com/garyburd/redigo/redis"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"github.com/thatisuday/commando"
	"os"
	"strings"
)

const JobNamespace = "supreme_robot"

var logger = loggo.GetLogger("main")

func main() {

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

			// Redis
			redisPool, err := initRedisPool()
			if err != nil {
				logger.Errorf("could not start redis: %s", err.Error())
				return
			}

			// Database
			modelClient, err := initModels(true)
			if err != nil {
				logger.Errorf("could not start redis: %s", err.Error())
				return
			}

			// Job Queue
			enqueuer := jobs.NewEnqueuer(JobNamespace, redisPool)

			// Webapp
			webServer, err := webapp.NewServer(redisPool, modelClient, enqueuer)
			if err != nil {
				logger.Errorf("could not start webapp: %s", err.Error())
				return
			}

			err = webServer.Run()
			if err != nil {
				logger.Errorf("Could not start webapp server %s", err.Error())
			}

		})

	commando.Parse(nil)

}

func initRedisPool() (*redis.Pool, error) {
	logger.Debugf("starting redis pool")

	// DB_ENGINE
	RedisAddress := os.Getenv("REDIS_ADDRESS")
	if RedisAddress == "" {
		return nil, errors.New("missing env var REDIS_ADDRESS")
	}

	return &redis.Pool{
		MaxActive: 10,
		MaxIdle:   10,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", RedisAddress)
		},
	}, nil
}

func initTwilio() (*twilio.Client, error) {
	// Twilio Stuff
	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	if twilioAccountSID == "" {
		return nil, errors.New("missing env var TWILIO_ACCOUNT_SID")
	}

	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if twilioAuthToken == "" {
		return nil, errors.New("missing env var TWILIO_TOKEN")
	}

	return twilio.NewClient(twilioAccountSID, twilioAuthToken), nil
}

func initModels(doMigration bool) (*models.Client, error) {
	// DB_ENGINE
	DBEngine := os.Getenv("DB_ENGINE")
	if DBEngine == "" {
		return nil, errors.New("missing env var DB_ENGINE")
	}

	return models.NewClient(DBEngine, doMigration)
}
