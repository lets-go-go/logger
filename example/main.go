package main

import (
	"os"
	"os/signal"

	"github.com/lets-go-go/logger"
	"github.com/robot4s/wechat/appconf"
)

func main() {
	config := logger.DefalutConfig()
	config.Level = logger.LEVEL(appconf.LogLevel)
	config.LogFileRollingType = logger.RollingDaily
	config.LogFileOutputDir = appconf.LogDir
	// config.LogFileName = "test"
	config.LogFileMaxCount = 5
	// config.LogFileMaxSize = 5
	// config.LogFileMaxSizeUnit = "MB"

	logger.Init(config)

	// for index := 0; index < 100000; index++ {
	logger.Trace("i am trace")
	logger.Debug("i am debug")
	logger.Warn("i am warning")
	logger.Error("i am error")
	logger.Fatal("i am fatal")
	// }

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan // wait for SIGINT

	// conf, _ := ioutil.ReadFile("conf.json")
}
