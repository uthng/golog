package golog_test

import (
	"github.com/uthng/golog"
	"io"
	"os"

	"github.com/fatih/color"
)

func ExampleStdout() {
	logger := golog.NewLogger(os.Stdout)
	logger.SetVerbosity(golog.DEBUG)

	logger.Debugln("This is debug log")
	logger.Infoln("This is info log")
	logger.Warnln("This is warn log")
	logger.Errorln("This is error log")
}

func ExampleMultiple() {
	file, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(-1)
	}

	multi := io.MultiWriter(file, os.Stdout)

	logger := golog.NewLogger(multi)
	logger.SetVerbosity(golog.INFO)

	red1 := color.New(color.FgRed)
	boldRed := red1.Add(color.Bold)
	boldRed.Println("This will print text in bold red.")

	logger.Debugln("This is debug log")
	logger.Infoln("This is info log")
	logger.Warnln("This is warn log")
	logger.Errorln("This is error log")
}
