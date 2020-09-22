package main

import (
	"github.com/creativeprojects/clog"
)

// setupLogger is configuring the default console handler
func setupLogger() {
	level := clog.LevelInfo
	logger := clog.NewFilteredConsoleLogger(level)
	clog.SetDefaultLogger(logger)
}
