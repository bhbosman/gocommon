package log_test

import (
	log2 "github.com/bhbosman/gocommon/log"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestLoggingLogLevel(t *testing.T) {
	tempMap := make(map[string]*log2.SubSystemLogger)
	fac := log2.NewFactory(log.New(os.Stdout, "", 0))
	subSystemLoggerFromFactory := fac.Create("ddd")
	tempMap["ddd"] = subSystemLoggerFromFactory
	var memoryClone log2.SubSystemLogger = *subSystemLoggerFromFactory
	fac.SetLogLevel("ddd", 12)
	assert.Equal(t, 12, tempMap["ddd"].GetLogLevel())
	assert.Equal(t, 12, subSystemLoggerFromFactory.GetLogLevel())
	assert.Equal(t, 12, memoryClone.GetLogLevel())

}
