package golog

import (
	"io"
	//"io/ioutil"
	"log"
	"os"
	//"fmt"

	"github.com/fatih/color"
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

var prefixes = map[int]string{
	ERROR: "ERROR: ",
	WARN:  "WARN: ",
	INFO:  "INFO: ",
	DEBUG: "DEBUG: ",
}

var colors = map[int][]color.Attribute{
	ERROR: []color.Attribute{color.FgRed},
	WARN:  []color.Attribute{color.FgYellow},
	INFO:  []color.Attribute{color.FgGreen},
	DEBUG: []color.Attribute{color.FgWhite},
}

// Logger is a wrapper of go log integrating log level
type Logger struct {
	loggers map[int]*log.Logger
	verbose int // if 0, no log
}

var defaultLogger *Logger

// Init a default logger with verbose = 3 and
// output for all loggers is stdout with different colors
func init() {
	defaultLogger = NewLogger(os.Stdout)
	defaultLogger.SetOutput(ERROR, os.Stderr)
}

// NewLogger returns a new instance logger
func NewLogger(w io.Writer) *Logger {
	color.NoColor = false
	logger := &Logger{}
	logger.verbose = 3

	logger.loggers = make(map[int]*log.Logger)
	for i := ERROR; i <= DEBUG; i++ {
		c := color.New(colors[i]...).Add(color.Bold).SprintFunc()
		logger.loggers[i] = log.New(w, c(prefixes[i]), log.Ldate|log.Ltime)
	}

	return logger
}

////////////////// INSTANCE LOGGER //////////////////////////////

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

// SetOutput sets output destination for a specific level
func (l *Logger) SetOutput(level int, w io.Writer) {
	l.loggers[level].SetOutput(w)
}

// SetFlags sets the output flags for a specific level
func (l *Logger) SetFlags(level int, flag int) {
	l.loggers[level].SetFlags(flag)
}

// GetFlags return the output flags of a specific level
func (l *Logger) GetFlags(level int) int {
	return l.loggers[level].Flags()
}

// Debug logs with debug level
func (l *Logger) Debug(v ...interface{}) {
	Print(l, DEBUG, v...)
}

// Debugf logs with debug level
func (l *Logger) Debugf(f string, v ...interface{}) {
	Printf(l, DEBUG, f, v...)
}

// Debugln logs with debug level
func (l *Logger) Debugln(v ...interface{}) {
	Println(l, DEBUG, v...)
}

// Info logs with info level
func (l *Logger) Info(v ...interface{}) {
	Print(l, INFO, v...)
}

// Infof logs with info level
func (l *Logger) Infof(f string, v ...interface{}) {
	Printf(l, INFO, f, v...)
}

// Infoln logs with info level
func (l *Logger) Infoln(v ...interface{}) {
	Println(l, INFO, v...)
}

// Warn logs with warn level
func (l *Logger) Warn(v ...interface{}) {
	Print(l, WARN, v...)
}

// Warnf logs with warn level
func (l *Logger) Warnf(f string, v ...interface{}) {
	Printf(l, WARN, f, v...)
}

// Warnln logs with warn level
func (l *Logger) Warnln(v ...interface{}) {
	Println(l, WARN, v...)
}

// Error logs with error level
func (l *Logger) Error(v ...interface{}) {
	Print(l, ERROR, v...)
}

// Errorf logs with error level
func (l *Logger) Errorf(f string, v ...interface{}) {
	Printf(l, ERROR, f, v...)
}

// Errorln logs with error level
func (l *Logger) Errorln(v ...interface{}) {
	Println(l, ERROR, v...)
}

//////////// DEFAULT LOGGER ////////////////////////////

// SetVerbosity sets log level. If verbose < NONE, it will be set to NONE.
// If verbose > DEBUG, it will be set to DEBUG
func SetVerbosity(v int) {
	if v < NONE {
		defaultLogger.verbose = NONE
	} else if v > DEBUG {
		defaultLogger.verbose = DEBUG
	} else {
		defaultLogger.verbose = v
	}
}

// GetVerbosity returns the current log level
func GetVerbosity() int {
	return defaultLogger.verbose
}

// SetOutput sets output destination for a specific level
func SetOutput(level int, w io.Writer) {
	defaultLogger.loggers[level].SetOutput(w)
}

// SetFlags sets the output flags for a specific level
func SetFlags(level int, flag int) {
	defaultLogger.loggers[level].SetFlags(flag)
}

// GetFlags return the output flags of a specific level
func GetFlags(level int) int {
	return defaultLogger.loggers[level].Flags()
}

// Debug logs with debug level
func Debug(v ...interface{}) {
	Print(defaultLogger, DEBUG, v...)
}

// Debugf logs with debug level
func Debugf(f string, v ...interface{}) {
	Printf(defaultLogger, DEBUG, f, v...)
}

// Debugln logs with debug level
func Debugln(v ...interface{}) {
	Println(defaultLogger, DEBUG, v...)
}

// Info logs with info level
func Info(v ...interface{}) {
	Print(defaultLogger, INFO, v...)
}

// Infof logs with info level
func Infof(f string, v ...interface{}) {
	Printf(defaultLogger, INFO, f, v...)
}

// Infoln logs with info level
func Infoln(v ...interface{}) {
	Println(defaultLogger, INFO, v...)
}

// Warn logs with warn level
func Warn(v ...interface{}) {
	Print(defaultLogger, WARN, v...)
}

// Warnf logs with warn level
func Warnf(f string, v ...interface{}) {
	Printf(defaultLogger, WARN, f, v...)
}

// Warnln logs with warn level
func Warnln(v ...interface{}) {
	Println(defaultLogger, WARN, v...)
}

// Error logs with error level
func Error(v ...interface{}) {
	Print(defaultLogger, ERROR, v...)
}

// Errorf logs with error level
func Errorf(f string, v ...interface{}) {
	Printf(defaultLogger, ERROR, f, v...)
}

// Errorln logs with error level
func Errorln(v ...interface{}) {
	Println(defaultLogger, ERROR, v...)
}

/////////////// INTERNAL FUNCTIONS /////////////////////

// Print wraps Print function of go log. It only prints
// log message if the level >= current verbose
func Print(l *Logger, level int, v ...interface{}) {
	if l.verbose >= level {
		c := color.New(colors[level]...).SprintFunc()
		l.loggers[level].Print(c(v...))
	}
}

// Printf wraps Printf function of go log. It only prints
// log message if the level >= current verbose
func Printf(l *Logger, level int, f string, v ...interface{}) {
	if l.verbose >= level {
		c := color.New(colors[level]...).SprintfFunc()
		l.loggers[level].Printf(c(f, v...))
	}
}

// Println wraps Println function of go log. It only prints
// log message if the level >= current verbose
func Println(l *Logger, level int, v ...interface{}) {
	if l.verbose >= level {
		c := color.New(colors[level]...).SprintlnFunc()
		l.loggers[level].Println(c(v...))
	}
}
