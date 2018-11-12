package golog_example

import (
	"github.com/uthng/golog"
)

func ExampleStdout() {
	logger := golog.NewLogger(&buf)
	logger.SetVerbosity(golog.DEBUG)

	logger.Debugln("This is debug log")
	logger.Infoln("This is info log")
	logger.Warnln("This is warn log")
	logger.Errorln("This is error log")
}
