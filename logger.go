package main

import (
	"log"

	"github.com/creativeprojects/clog"
)

// setupLogger is configuring the default logger output (console or file)
// and returns a function to close the logs at the end
func setupLogger(flags *Flags) func() {
	var err error
	var handler clog.Handler

	// by default, nothing to close at the end
	close := func() {}

	if flags.outputFilename != "" {
		handler, err = clog.NewFileHandler(flags.outputFilename, "", log.LstdFlags)
		if err != nil {
			// open the console as a backup
			handler = clog.NewConsoleHandler("", 0)
			// and pushes a warning manually (there should be a better way of doing this?)
			handler.LogEntry(clog.LogEntry{
				Level:  clog.LevelWarning,
				Values: []interface{}{"cannot open output file: logging to the console instead"},
			})
		} else {
			// will have to close the file at the end
			close = handler.(*clog.FileHandler).Close
		}
	} else {
		handler = clog.NewConsoleHandler("", 0)
	}

	level := clog.LevelInfo
	if flags.verbose {
		level = clog.LevelDebug
	}
	logger := clog.NewLogger(clog.NewLevelFilter(level, handler))
	clog.SetDefaultLogger(logger)

	return close
}
