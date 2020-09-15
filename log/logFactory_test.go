package log

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestLoggingLogLevel(t *testing.T) {
	tempMap := make(map[string]*SubSystemLogger)
	fac := NewLogFactory(log.New(os.Stdout, "", 0))
	subSystemLoggerFromFactory := fac.Create("ddd")
	tempMap["ddd"] = subSystemLoggerFromFactory
	var memoryClone SubSystemLogger = *subSystemLoggerFromFactory
	fac.SetLogLevel("ddd", 12)
	assert.Equal(t, 12, tempMap["ddd"].GetLogLevel())
	assert.Equal(t, 12, subSystemLoggerFromFactory.GetLogLevel())
	assert.Equal(t, 12, memoryClone.GetLogLevel())

}
