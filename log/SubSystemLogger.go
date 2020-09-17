package log

import (
	"fmt"
	"log"
)

type SubSystemLogger struct {
	l           *log.Logger
	initialized bool
	registerCb  func(string) int
	name        string
	renameCb    func(s string) *log.Logger
}

func NewSubSystemLogger(name string, l *log.Logger, registerCb func(string) int, renameCb func(s string) *log.Logger) *SubSystemLogger {
	return &SubSystemLogger{
		l:           l,
		initialized: l != nil,
		registerCb:  registerCb,
		name:        name,
		renameCb:    renameCb,
	}
}

func (ssl SubSystemLogger) Printf(s string, v ...interface{}) {
	ssl.LogWithLevel(0, func(logger *log.Logger) {
		logger.Printf(s, v...)
	})
}

func (ssl SubSystemLogger) Name() string {
	return ssl.name
}

func (l SubSystemLogger) GetLogLevel() int {
	if l.registerCb != nil {
		return l.registerCb(l.name)
	}
	return 0
}

func (ssl SubSystemLogger) LogWithLevel(level int, cb func(logger *log.Logger)) {
	l := ssl.GetLogLevel()
	if level <= l {
		if cb != nil {
			cb(ssl.l)
		}
	}
}

func (ssl SubSystemLogger) Error(err error) error {
	if err != nil {
		if ssl.initialized {
			ssl.l.Printf("Error: %v", err.Error())
		}
	}
	return err
}

func (ssl SubSystemLogger) ErrorWithDescription(description string, err error) error {
	if err != nil {
		if ssl.initialized {
			ssl.l.Printf("Error Message: %v, err:  %v", description, err.Error())
		}
	}
	return err
}

func (ssl *SubSystemLogger) NameChange(newName string) {
	ssl.LogWithLevel(0, func(logger *log.Logger) {
		logger.Printf(fmt.Sprintf("Change name from %v to %v", ssl.name, newName))
	})
	ssl.l = ssl.renameCb(newName)
}
