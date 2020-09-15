package log

import (
	"fmt"
	"log"
)

type LogWriter struct {
	l      *log.Logger
	prefix string
}

func NewLogWriter(l *log.Logger, prefix string) *LogWriter {
	return &LogWriter{
		l:      l,
		prefix: prefix,
	}
}

func (lw LogWriter) Write(p []byte) (n int, err error) {
	s := string(p)
	lw.l.Printf("[%v] %v", lw.prefix, s)
	return len(p), nil
}

type LogFactory struct {
	l      *log.Logger
	levels map[string]int
}

func NewLogFactory(l *log.Logger) *LogFactory {
	return &LogFactory{
		l:      l,
		levels: make(map[string]int),
	}
}

type SubSystemLogger struct {
	l           *log.Logger
	initialized bool
	registerCb  func(string) int
	name        string
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

func NewSubSystemLogger(name string, l *log.Logger, registerCb func(string) int) *SubSystemLogger {
	return &SubSystemLogger{
		l:           l,
		initialized: l != nil,
		registerCb:  registerCb,
		name:        name,
	}
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

func (ssl SubSystemLogger) Sub(s string) SubSystemLogger {
	return *NewSubSystemLogger(fmt.Sprintf("%v-%v", ssl.name, s), ssl.l, ssl.registerCb)
}

func (lf LogFactory) Create(subSystem string) *SubSystemLogger {
	return NewSubSystemLogger(
		subSystem,
		log.New(NewLogWriter(lf.l, subSystem), "", 0),
		func(s string) int {
			if v, ok := lf.levels[s]; ok {
				return v
			}
			return 0
		})
}

func (lf *LogFactory) SetLogLevel(s string, i int) {
	lf.levels[s] = i
}
