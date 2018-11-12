package golog

import (
	"io"
	//"io/ioutil"
	"log"
	//"os"
	//"fmt"
)

const (
	// NONE = 0
	NONE = iota
	// ERROR = 1
	ERROR
	// WARN = 2
	WARN
	// INFO = 3
	INFO
	// DEBUG = 4
	DEBUG
)

var levels = map[int]string{
	ERROR: "ERROR: ",
	WARN:  "WARN: ",
	INFO:  "INFO: ",
	DEBUG: "DEBUG: ",
}

// Logger is a wrapper of go log integrating log level
type Logger struct {
	logger *log.Logger

	verbose int // if 0, no log
}

// NewLogger initializes a logger with a io writer for all log levels.
func NewLogger(w io.Writer) *Logger {
	return &Logger{
		logger:  log.New(w, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		verbose: 3,
	}
}

// SetVerbosity sets log level. If verbose < NONE, it will be set to NONE.
// If verbose > DEBUG, it will be set to DEBUG
func (l *Logger) SetVerbosity(v int) {
	if v < NONE {
		l.verbose = NONE
	} else if v > DEBUG {
		l.verbose = DEBUG
	} else {
		l.verbose = v
	}
}

// GetVerbosity returns the current log level
func (l *Logger) GetVerbosity() int {
	return l.verbose
}

// SetOutput sets output destination for the logger
func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

// SetFlags sets the output flags for the logger
func (l *Logger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

// GetFlags return the output flags of the logger
func (l *Logger) GetFlags() int {
	return l.logger.Flags()
}

// Debug logs with debug level
func (l *Logger) Debug(v ...interface{}) {
	l.Print(DEBUG, v...)
}

// Debugf logs with debug level
func (l *Logger) Debugf(f string, v ...interface{}) {
	l.Printf(DEBUG, f, v...)
}

// Debugln logs with debug level
func (l *Logger) Debugln(v ...interface{}) {
	l.Println(DEBUG, v...)
}

// Info logs with info level
func (l *Logger) Info(v ...interface{}) {
	l.Print(INFO, v...)
}

// Infof logs with info level
func (l *Logger) Infof(f string, v ...interface{}) {
	l.Printf(INFO, f, v...)
}

// Infoln logs with info level
func (l *Logger) Infoln(v ...interface{}) {
	l.Println(INFO, v...)
}

// Warn logs with warn level
func (l *Logger) Warn(v ...interface{}) {
	l.Print(WARN, v...)
}

// Warnf logs with warn level
func (l *Logger) Warnf(f string, v ...interface{}) {
	l.Printf(WARN, f, v...)
}

// Warnln logs with warn level
func (l *Logger) Warnln(v ...interface{}) {
	l.Println(WARN, v...)
}

// Error logs with error level
func (l *Logger) Error(v ...interface{}) {
	l.Print(ERROR, v...)
}

// Errorf logs with error level
func (l *Logger) Errorf(f string, v ...interface{}) {
	l.Printf(ERROR, f, v...)
}

// Errorln logs with error level
func (l *Logger) Errorln(v ...interface{}) {
	l.Println(ERROR, v...)
}

// Print wraps Print function of go log. It only prints
// log message if the level >= current verbose
func (l *Logger) Print(level int, v ...interface{}) {
	if l.verbose >= level {
		l.logger.SetPrefix(levels[level])
		l.logger.Print(v...)
	}
}

// Printf wraps Printf function of go log. It only prints
// log message if the level >= current verbose
func (l *Logger) Printf(level int, f string, v ...interface{}) {
	if l.verbose >= level {
		l.logger.SetPrefix(levels[level])
		l.logger.Printf(f, v...)
	}
}

// Println wraps Println function of go log. It only prints
// log message if the level >= current verbose
func (l *Logger) Println(level int, v ...interface{}) {
	if l.verbose >= level {
		l.logger.SetPrefix(levels[level])
		l.logger.Println(v...)
	}
}
