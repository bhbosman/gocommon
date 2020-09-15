package log

import "log"

type LoggerWrapper struct {
	ApplicationLogger *log.Logger
}

func (self *LoggerWrapper) Printf(format string, v ...interface{}) {
	self.ApplicationLogger.Printf(format, v...)
}
