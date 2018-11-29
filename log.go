package golog

import (
	"io"
	//"io/ioutil"
	"fmt"
	//"log"
	"os"
	"path"
	//"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cast"
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

const (
	// PRINT = 0
	PRINT = iota
	// PRINTF = 1
	PRINTF
	// PRINTLN = 2
	PRINTLN
	// PRINTW = 3
	PRINTW
)

var prefixes = map[int]string{
	ERROR: "ERROR:",
	WARN:  "WARN:",
	INFO:  "INFO:",
	DEBUG: "DEBUG:",
}

var colors = map[int][]color.Attribute{
	ERROR: []color.Attribute{color.FgRed},
	WARN:  []color.Attribute{color.FgYellow},
	INFO:  []color.Attribute{color.FgGreen},
	DEBUG: []color.Attribute{color.FgWhite},
}

var mutex sync.Mutex
var wg sync.WaitGroup

type level struct {
	output      io.Writer
	color       bool
	colorPrefix *color.Color
	colorText   *color.Color
}

// Logger is a wrapper of go log integrating log level
type Logger struct {
	levels  map[int]*level
	verbose int // if 0, no log
	caller  bool
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
	logger.caller = false

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
			output:      w,
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
		l.levels[i].output = w
	}
}

// SetLevelOutput sets output destination for a specific level
func (l *Logger) SetLevelOutput(level int, w io.Writer) {
	l.levels[level].output = w
}

// EnableCaller enables/disables caller infos for all log levels
func (l *Logger) EnableCaller(enabled bool) {
	l.caller = enabled
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
	//l.levels[level].SetPrefix(cf.SprintFunc()(prefixes[level]))
}

// DisableLevelColor enables color for a specific level
func (l *Logger) DisableLevelColor(level int) {
	l.levels[level].color = false
	cf := l.levels[level].colorPrefix
	cf.DisableColor()
	//l.levels[level].SetPrefix(cf.SprintFunc()(prefixes[level]))
}

// Debug logs with debug level
func (l *Logger) Debug(v ...interface{}) {
	Log(PRINT, l, DEBUG, "", v...)
}

// Debugf logs with debug level
func (l *Logger) Debugf(f string, v ...interface{}) {
	Log(PRINTF, l, DEBUG, f, v...)
}

// Debugln logs with debug level
func (l *Logger) Debugln(v ...interface{}) {
	Log(PRINTLN, l, DEBUG, "", v...)
}

// Debugw logs with debug level with structured log format
func (l *Logger) Debugw(msg string, v ...interface{}) {
	Log(PRINTW, l, DEBUG, msg, v...)
}

// Info logs with info level
func (l *Logger) Info(v ...interface{}) {
	Log(PRINT, l, INFO, "", v...)
}

// Infof logs with info level
func (l *Logger) Infof(f string, v ...interface{}) {
	Log(PRINTF, l, INFO, f, v...)
}

// Infoln logs with info level
func (l *Logger) Infoln(v ...interface{}) {
	Log(PRINTLN, l, INFO, "", v...)
}

// Warn logs with warn level
func (l *Logger) Warn(v ...interface{}) {
	Log(PRINT, l, WARN, "", v...)
}

// Warnf logs with warn level
func (l *Logger) Warnf(f string, v ...interface{}) {
	Log(PRINTF, l, WARN, f, v...)
}

// Warnln logs with warn level
func (l *Logger) Warnln(v ...interface{}) {
	Log(PRINTLN, l, WARN, "", v...)
}

// Error logs with error level
func (l *Logger) Error(v ...interface{}) {
	Log(PRINT, l, ERROR, "", v...)
}

// Errorf logs with error level
func (l *Logger) Errorf(f string, v ...interface{}) {
	Log(PRINTF, l, ERROR, f, v...)
}

// Errorln logs with error level
func (l *Logger) Errorln(v ...interface{}) {
	Log(PRINTLN, l, ERROR, "", v...)
}

// Fatal logs with Print() followed by os.Exit(1)
func (l *Logger) Fatal(v ...interface{}) {
	Log(PRINT, l, FATAL, "", v...)
	os.Exit(1)
}

// Fatalf logs with Printf() followed by os.Exit(1)
func (l *Logger) Fatalf(f string, v ...interface{}) {
	Log(PRINTF, l, FATAL, f, v...)
	os.Exit(1)
}

// Fatalln logs with Println() followed by os.Exit(1)
func (l *Logger) Fatalln(v ...interface{}) {
	Log(PRINTLN, l, FATAL, "", v...)
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
		defaultLogger.levels[i].output = w
	}
}

// SetLevelOutput sets output destination for a specific level
func SetLevelOutput(level int, w io.Writer) {
	defaultLogger.levels[level].output = w
}

// EnableCaller enables/disables caller infos for all log levels
func EnableCaller(enabled bool) {
	defaultLogger.caller = enabled
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

}

// DisableLevelColor enables color for a specific level
func DisableLevelColor(level int) {
	defaultLogger.levels[level].color = false
	cf := defaultLogger.levels[level].colorPrefix
	cf.DisableColor()
}

// Debug logs with debug level
func Debug(v ...interface{}) {
	Log(PRINT, defaultLogger, DEBUG, "", v...)
}

// Debugf logs with debug level
func Debugf(f string, v ...interface{}) {
	Log(PRINTF, defaultLogger, DEBUG, f, v...)
}

// Debugln logs with debug level
func Debugln(v ...interface{}) {
	Log(PRINTLN, defaultLogger, DEBUG, "", v...)
}

// Info logs with info level
func Info(v ...interface{}) {
	Log(PRINT, defaultLogger, INFO, "", v...)
}

// Infof logs with info level
func Infof(f string, v ...interface{}) {
	Log(PRINTF, defaultLogger, INFO, f, v...)
}

// Infoln logs with info level
func Infoln(v ...interface{}) {
	Log(PRINTLN, defaultLogger, INFO, "", v...)
}

// Warn logs with warn level
func Warn(v ...interface{}) {
	Log(PRINT, defaultLogger, WARN, "", v...)
}

// Warnf logs with warn level
func Warnf(f string, v ...interface{}) {
	Log(PRINTF, defaultLogger, WARN, f, v...)
}

// Warnln logs with warn level
func Warnln(v ...interface{}) {
	Log(PRINTLN, defaultLogger, WARN, "", v...)
}

// Error logs with error level
func Error(v ...interface{}) {
	Log(PRINT, defaultLogger, ERROR, "", v...)
}

// Errorf logs with error level
func Errorf(f string, v ...interface{}) {
	Log(PRINTF, defaultLogger, ERROR, f, v...)
}

// Errorln logs with error level
func Errorln(v ...interface{}) {
	Log(PRINTLN, defaultLogger, ERROR, "", v...)
}

// Fatal logs with Print() followed by os.Exit(1)
func Fatal(v ...interface{}) {
	Log(PRINT, defaultLogger, FATAL, "", v...)
	os.Exit(1)
}

// Fatalf logs with Printf() followed by os.Exit(1)
func Fatalf(f string, v ...interface{}) {
	Log(PRINTF, defaultLogger, FATAL, f, v...)
	os.Exit(1)
}

// Fatalln logs with Println() followed by os.Exit(1)
func Fatalln(v ...interface{}) {
	Log(PRINTLN, defaultLogger, FATAL, "", v...)
	os.Exit(1)
}

/////////////// INTERNAL FUNCTIONS /////////////////////

// Log wraps print function but using goroutine and waitgroup
// to have a synchronization of logs.
func Log(p int, l *Logger, level int, f string, v ...interface{}) {
	wg.Add(1)
	caller := getInfoCaller()
	go func(c string) {
		defer wg.Done()
		printMsg(p, l, level, caller, f, v...)
	}(caller)
	wg.Wait()
}

func printMsg(p int, l *Logger, level int, caller string, f string, v ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	var prefix string

	if l.verbose >= level {
		ct := l.levels[level].colorText
		cf := l.levels[level].colorPrefix

		if l.levels[level].color {
			ct.EnableColor()
			cf.EnableColor()
		} else {
			ct.DisableColor()
			cf.DisableColor()
		}

		if l.caller {
			prefix = fmt.Sprintf("%s %s %s", getTimeNow(), caller, cf.SprintFunc()(prefixes[level]))
		} else {
			prefix = fmt.Sprintf("%s %s", getTimeNow(), cf.SprintFunc()(prefixes[level]))
		}

		switch p {
		case PRINT:
			fmt.Fprint(l.levels[level].output, prefix, " ", ct.SprintFunc()(v...))
		case PRINTF:
			fmt.Fprintf(l.levels[level].output, "%s %s", prefix, ct.SprintfFunc()(f, v...))
		case PRINTLN:
			fmt.Fprintln(l.levels[level].output, prefix, ct.SprintlnFunc()(v...))
		case PRINTW:
			printw(l.levels[level].output, prefix, ct, f, v...)
		default:
			fmt.Fprintln(l.levels[level].output, prefix, ct.SprintlnFunc()(v...))
		}
	}
}

func getInfoCaller() string {
	if pc, file, line, ok := runtime.Caller(3); ok {
		fn := runtime.FuncForPC(pc).Name()
		arr := strings.Split(path.Base(fn), ".")
		str := fmt.Sprintf("%s:%d:%s", path.Base(file), line, arr[len(arr)-1])
		return str
	}

	return ""
}

func getTimeNow() string {
	return time.Now().Format(time.RFC3339Nano)
}

func printw(output io.Writer, prefix string, ct *color.Color, msg string, keyvals ...interface{}) {
	var pairs []interface{}
	var format string
	kv := keyvals

	pairs = append(pairs, prefix, ct.SprintFunc()(msg))
	format += "%s %s "

	if len(kv)%2 != 0 {
		kv = append(kv, "missing")
	}

	for i := 0; i < len(kv); i += 2 {
		// cast 1st elem = key to string
		k := cast.ToString(kv[i])
		if k == "" {
			k = "missing"
		}
		k = ct.SprintFunc()(k)

		// cast 2nd elem = value
		v := kv[i+1]
		pair := ""
		//kind := reflect.ValueOf(v).Kind()

		//if kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map || kind == reflect.Struct || kind == reflect.Ptr {
		//pair = fmt.Sprintf("%s=%+v", k, v)
		//} else {
		//s := cast.ToString(v)
		//if strings.Contains(s, " ") {
		//pair = fmt.Sprintf("%s=\"%s\"", k, cast.ToString(v))
		//} else {
		//pair = fmt.Sprintf("%s=%s", k, cast.ToString(v))
		//}
		//}
		pair = fmt.Sprintf("%s=%+v", k, v)
		pairs = append(pairs, pair)
		if i != len(kv)-2 {
			format += "%s "
		} else {
			format += "%s\n"
		}
	}

	fmt.Fprintf(output, format, pairs...)
}
