package golog

import (
	"io"
	//"io/ioutil"
	"log"
	//"os"
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

var levels map[int]string
var colors map[int][]color.Attribute

// Logger is a wrapper of go log integrating log level
type Logger struct {
	logger *log.Logger

	verbose int // if 0, no log
}

// NewLogger initializes a logger with a io writer for all log levels.
func NewLogger(w io.Writer) *Logger {
	color.NoColor = false

	colors = make(map[int][]color.Attribute)
	colors[ERROR] = []color.Attribute{color.FgRed}
	colors[WARN] = []color.Attribute{color.FgYellow}
	colors[INFO] = []color.Attribute{color.FgGreen}
	colors[DEBUG] = []color.Attribute{color.FgWhite}

	red := color.New(colors[ERROR]...).Add(color.Bold).SprintFunc()
	yellow := color.New(colors[WARN]...).Add(color.Bold).SprintFunc()
	green := color.New(colors[INFO]...).Add(color.Bold).SprintFunc()
	white := color.New(colors[DEBUG]...).Add(color.Bold).SprintFunc()

	levels = make(map[int]string)
	levels[ERROR] = red("ERROR: ")
	levels[WARN] = yellow("WARN: ")
	levels[INFO] = green("INFO: ")
	levels[DEBUG] = white("DEBUG: ")

	return &Logger{
		logger:  log.New(w, levels[INFO], log.Ldate|log.Ltime),
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
		c := color.New(colors[level]...).SprintFunc()
		l.logger.Print(c(v...))
	}
}

// Printf wraps Printf function of go log. It only prints
// log message if the level >= current verbose
func (l *Logger) Printf(level int, f string, v ...interface{}) {
	if l.verbose >= level {
		l.logger.SetPrefix(levels[level])
		c := color.New(colors[level]...).SprintfFunc()
		l.logger.Printf(c(f, v...))
	}
}

// Println wraps Println function of go log. It only prints
// log message if the level >= current verbose
func (l *Logger) Println(level int, v ...interface{}) {
	if l.verbose >= level {
		l.logger.SetPrefix(levels[level])
		c := color.New(colors[level]...).SprintlnFunc()
		l.logger.Println(c(v...))
	}
}
