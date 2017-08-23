package main

import (
	"github.com/lets-go-go/logger"
)

func main() {
	config := logger.DefalutConfig()
	config.Level = logger.DEBUG
	config.LogFileRollingType = logger.RollingSize
	config.LogFileMaxSize = 5
	config.LogFileMaxSizeUnit = "MB"

	logger.Init(config)

	for index := 0; index < 100000; index++ {
		logger.Trace("i am trace")
		logger.Debug("i am debug")
		logger.Warn("i am warning")
		logger.Error("i am error")
		logger.Fatal("i am fatal")
	}

	// conf, _ := ioutil.ReadFile("conf.json")
}
