package golog

import (
	"io"
	//"io/ioutil"
	"log"
	"os"
	//"fmt"
	"sync"

	"github.com/fatih/color"
)

const (
	// NONE = 0
	NONE = iota
	// FATAL = 1
	FATAL
	// ERROR = 2
	ERROR
	// WARN = 3
	WARN
	// INFO = 4
	INFO
	// DEBUG = 5
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

var mutex sync.Mutex

type level struct {
	*log.Logger
	color       bool
	colorPrefix *color.Color
	colorText   *color.Color
}

// Logger is a wrapper of go log integrating log level
type Logger struct {
	levels  map[int]*level
	verbose int // if 0, no log
}

var defaultLogger *Logger

// Init a default logger with verbose = 3 and
// output for all levels is stdout with different colors
func init() {
	defaultLogger = NewLogger()
}

// NewLogger returns a new instance logger
// By default, it uses stderr for error and stdout for other levels
func NewLogger() *Logger {
	color.NoColor = false
	logger := &Logger{}
	logger.verbose = 4

	logger.levels = make(map[int]*level)
	for i := FATAL; i <= DEBUG; i++ {
		colorPrefix := color.New(colors[i]...).Add(color.Bold)
		colorText := color.New(colors[i]...)
		w := os.Stdout
		if i == FATAL || i == ERROR {
			w = os.Stderr
		}
		logger.levels[i] = &level{
			color:       true,
			colorPrefix: colorPrefix,
			colorText:   colorText,
			Logger:      log.New(w, colorPrefix.SprintFunc()(prefixes[i]), log.Ldate|log.Ltime),
		}
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
func (l *Logger) SetOutput(w io.Writer) {
	for i := FATAL; i <= DEBUG; i++ {
		l.levels[i].SetOutput(w)
	}
}

// SetLevelOutput sets output destination for a specific level
func (l *Logger) SetLevelOutput(level int, w io.Writer) {
	l.levels[level].SetOutput(w)
}

// SetFlags sets the same output flags for all levels
func (l *Logger) SetFlags(flag int) {
	for i := FATAL; i <= DEBUG; i++ {
		l.levels[i].SetFlags(flag)
	}
}

// SetLevelFlags sets the output flags for a specific level
func (l *Logger) SetLevelFlags(level int, flag int) {
	l.levels[level].SetFlags(flag)
}

// GetFlags returns the output flags of a specific level
func (l *Logger) GetFlags(level int) int {
	return l.levels[level].Flags()
}

// EnableColor enables color for all log levels
func (l *Logger) EnableColor() {
	for i := FATAL; i <= DEBUG; i++ {
		l.EnableLevelColor(i)
	}
}

// DisableColor disables color for all log levels
func (l *Logger) DisableColor() {
	for i := FATAL; i <= DEBUG; i++ {
		l.DisableLevelColor(i)
	}
}

// EnableLevelColor enables color for a specific level
func (l *Logger) EnableLevelColor(level int) {
	l.levels[level].color = true
	cf := l.levels[level].colorPrefix
	cf.EnableColor()
	l.levels[level].SetPrefix(cf.SprintFunc()(prefixes[level]))
}

// DisableLevelColor enables color for a specific level
func (l *Logger) DisableLevelColor(level int) {
	l.levels[level].color = false
	cf := l.levels[level].colorPrefix
	cf.DisableColor()
	l.levels[level].SetPrefix(cf.SprintFunc()(prefixes[level]))
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

// Fatal logs with Print() followed by os.Exit(1)
func (l *Logger) Fatal(v ...interface{}) {
	Print(l, FATAL, v...)
	os.Exit(1)
}

// Fatalf logs with Printf() followed by os.Exit(1)
func (l *Logger) Fatalf(f string, v ...interface{}) {
	Printf(l, FATAL, f, v...)
	os.Exit(1)
}

// Fatalln logs with Println() followed by os.Exit(1)
func (l *Logger) Fatalln(v ...interface{}) {
	Println(l, FATAL, v...)
	os.Exit(1)
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
func SetOutput(w io.Writer) {
	for i := FATAL; i <= DEBUG; i++ {
		defaultLogger.levels[i].SetOutput(w)
	}
}

// SetLevelOutput sets output destination for a specific level
func SetLevelOutput(level int, w io.Writer) {
	defaultLogger.levels[level].SetOutput(w)
}

// SetFlags sets the same output flags for all levels
func SetFlags(flag int) {
	for i := FATAL; i <= DEBUG; i++ {
		defaultLogger.levels[i].SetFlags(flag)
	}
}

// SetLevelFlags sets the output flags for a specific level
func SetLevelFlags(level int, flag int) {
	defaultLogger.levels[level].SetFlags(flag)
}

// GetFlags return the output flags of a specific level
func GetFlags(level int) int {
	return defaultLogger.levels[level].Flags()
}

// EnableColor enables color for all log levels
func EnableColor() {
	for i := FATAL; i <= DEBUG; i++ {
		defaultLogger.EnableLevelColor(i)
	}
}

// DisableColor disables color for all log levels
func DisableColor() {
	for i := FATAL; i <= DEBUG; i++ {
		defaultLogger.DisableLevelColor(i)
	}
}

// EnableLevelColor enables color for a specific level
func EnableLevelColor(level int) {
	defaultLogger.levels[level].color = true
	cf := defaultLogger.levels[level].colorPrefix
	cf.EnableColor()
	defaultLogger.levels[level].SetPrefix(cf.SprintFunc()(prefixes[level]))

}

// DisableLevelColor enables color for a specific level
func DisableLevelColor(level int) {
	defaultLogger.levels[level].color = false
	cf := defaultLogger.levels[level].colorPrefix
	cf.DisableColor()
	defaultLogger.levels[level].SetPrefix(cf.SprintFunc()(prefixes[level]))
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

// Fatal logs with Print() followed by os.Exit(1)
func Fatal(v ...interface{}) {
	Print(defaultLogger, FATAL, v...)
	os.Exit(1)
}

// Fatalf logs with Printf() followed by os.Exit(1)
func Fatalf(f string, v ...interface{}) {
	Printf(defaultLogger, FATAL, f, v...)
	os.Exit(1)
}

// Fatalln logs with Println() followed by os.Exit(1)
func Fatalln(v ...interface{}) {
	Println(defaultLogger, FATAL, v...)
	os.Exit(1)
}

/////////////// INTERNAL FUNCTIONS /////////////////////

// Print wraps Print function of go log. It only prints
// log message if the level >= current verbose
func Print(l *Logger, level int, v ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	if l.verbose >= level {
		ct := l.levels[level].colorText
		if l.levels[level].color {
			ct.EnableColor()
		} else {
			ct.DisableColor()
		}

		l.levels[level].Print(ct.SprintFunc()(v...))
	}
}

// Printf wraps Printf function of go log. It only prints
// log message if the level >= current verbose
func Printf(l *Logger, level int, f string, v ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	if l.verbose >= level {
		ct := l.levels[level].colorText
		if l.levels[level].color {
			ct.EnableColor()
		} else {
			ct.DisableColor()
		}
		l.levels[level].Printf(ct.SprintfFunc()(f, v...))
	}
}

// Println wraps Println function of go log. It only prints
// log message if the level >= current verbose
func Println(l *Logger, level int, v ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	if l.verbose >= level {
		ct := l.levels[level].colorText
		if l.levels[level].color {
			ct.EnableColor()
		} else {
			ct.DisableColor()
		}
		l.levels[level].Println(ct.SprintlnFunc()(v...))
	}
}
