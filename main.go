//go:generate pkger
package main

import (
	"errors"
	"fmt"
	"github.com/Hive-Gay/go-hivelib/clients/v1/quotes"
	"github.com/Hive-Gay/supreme-robot/config"
	"github.com/Hive-Gay/supreme-robot/database"
	"github.com/Hive-Gay/supreme-robot/redis"
	"github.com/Hive-Gay/supreme-robot/webapp"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logger = loggo.GetLogger("main")

func main() {
	c := config.CollectConfig()

	// start logger
	err := loggo.ConfigureLoggers(c.LoggerConfig)
	if err != nil {
		log.Fatalf("error configuring logger: %s", err.Error())
		return
	}
	_, err = loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	if err != nil {
		log.Fatalf("error configuring color logger: %s", err.Error())
		return
	}

	logger := loggo.GetLogger("main")
	logger.Infof("starting main process")

	// create redis client
	rc, err := redis.NewClient(c)
	if err != nil {
		logger.Errorf("new db client: %s", err.Error())
		return
	}

	// create database client
	db, err := database.NewClient(c)
	if err != nil {
		logger.Errorf("new db client: %s", err.Error())
		return
	}

	// create quotes client
	qc, err := quotes.NewClient()
	if err != nil {
		logger.Errorf("new quotes client: %s", err.Error())
		return
	}

	// create web server
	ws, err := webapp.NewServer(c, rc, db, qc)
	if err != nil {
		logger.Errorf("new db client: %s", err.Error())
		return
	}

	// ** start application **
	errChan := make(chan error)

	// start web server
	logger.Debugf("starting web app")
	go func(errChan chan error) {
		err := ws.ListenAndServe()
		if err != nil {
			errChan <- errors.New(fmt.Sprintf("webapp: %s", err.Error()))
		}
	}(errChan)
	defer ws.Close()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-nch:
		logger.Infof("got sig: %s", sig)
	case err := <-errChan:
		logger.Criticalf(err.Error())
	}

	logger.Infof("main process done")
}
