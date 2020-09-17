package log

import (
	"log"
)

type Writer struct {
	l      *log.Logger
	prefix string
}

func NewLogWriter(l *log.Logger, prefix string) *Writer {
	return &Writer{
		l:      l,
		prefix: prefix,
	}
}

func (lw Writer) Write(p []byte) (n int, err error) {
	s := string(p)
	lw.l.Printf("[%v] %v", lw.prefix, s)
	return len(p), nil
}

type Factory struct {
	l      *log.Logger
	levels map[string]int
}

func NewFactory(l *log.Logger) *Factory {
	return &Factory{
		l:      l,
		levels: make(map[string]int),
	}
}

func (lf Factory) Create(subSystem string) *SubSystemLogger {
	return NewSubSystemLogger(
		subSystem,
		log.New(NewLogWriter(lf.l, subSystem), "", 0),
		func(s string) int {
			if v, ok := lf.levels[s]; ok {
				return v
			}
			return 0
		},
		func(s string) *log.Logger {
			return log.New(NewLogWriter(lf.l, s), "", 0)
		})
}

func (lf *Factory) SetLogLevel(s string, i int) {
	lf.levels[s] = i
}
