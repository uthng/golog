package golog

import (
	"io"
	//"io/ioutil"
	"fmt"
	//"log"
	"os"
	"path"
	"reflect"
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

const (
	// FTIMESTAMP enables timestamp field in message log
	FTIMESTAMP = 1 << iota
	// FCALLER enables caller field in message log
	FCALLER
	// FFULLSTRUCTUREDLOG enables structured log for all fields in message log
	FFULLSTRUCTUREDLOG
)

var prefixes = map[int]string{
	FATAL: "FATAL",
	ERROR: "ERROR",
	WARN:  "WARN",
	INFO:  "INFO",
	DEBUG: "DEBUG",
}

var colors = map[int][]color.Attribute{
	FATAL: []color.Attribute{color.FgRed},
	ERROR: []color.Attribute{color.FgRed},
	WARN:  []color.Attribute{color.FgYellow},
	INFO:  []color.Attribute{color.FgGreen},
	DEBUG: []color.Attribute{color.FgWhite},
}

var mutex Mutex
var wg sync.WaitGroup

type level struct {
	output      io.Writer
	color       bool
	colorPrefix *color.Color
	colorText   *color.Color
}

// Logger is a wrapper of go log integrating log level
type Logger struct {
	levels     map[int]*level
	verbose    int // if 0, no log
	flag       int
	timeFormat string
	logFormat  bool

	handlers []Handler
}

// Handler defines an interface for golog handler
type Handler interface {
	PrintMsg(p int, l *Logger, level int, fields Fields) error
}

// Field defines a log field in case of structured log
type Field struct {
	Key   string
	Value string
}

// Fields stores all log fields: prefix, static and user log
type Fields struct {
	Prefix []*Field
	Log    []*Field
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
	logger.flag = 0 // no flag
	logger.timeFormat = time.RFC3339
	logger.logFormat = true

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
	mutex.LockOnce()
	defer mutex.UnlockOnce()

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
	mutex.LockOnce()
	defer mutex.UnlockOnce()

	for i := FATAL; i <= DEBUG; i++ {
		l.levels[i].output = w
	}
}

// SetLevelOutput sets output destination for a specific level
func (l *Logger) SetLevelOutput(level int, w io.Writer) {
	mutex.LockOnce()
	defer mutex.UnlockOnce()

	l.levels[level].output = w
}

// SetFlags sets flags for message log output
func (l *Logger) SetFlags(flag int) {
	mutex.LockOnce()
	defer mutex.UnlockOnce()

	l.flag = flag
}

// GetFlags gets flags for message log output
func (l *Logger) GetFlags() int {
	return l.flag
}

// SetTimeFormat sets timestamp with the given format
func (l *Logger) SetTimeFormat(format string) {
	mutex.LockOnce()
	defer mutex.UnlockOnce()

	l.timeFormat = format
}

// EnableColor enables color for all log levels
func (l *Logger) EnableColor() {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	for i := FATAL; i <= DEBUG; i++ {
		l.EnableLevelColor(i)
	}
}

// DisableColor disables color for all log levels
func (l *Logger) DisableColor() {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	for i := FATAL; i <= DEBUG; i++ {
		l.DisableLevelColor(i)
	}
}

// EnableLevelColor enables color for a specific level
func (l *Logger) EnableLevelColor(level int) {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	l.levels[level].color = true
	cf := l.levels[level].colorPrefix
	cf.EnableColor()
	//l.levels[level].SetPrefix(cf.SprintFunc()(prefixes[level]))
}

// DisableLevelColor enables color for a specific level
func (l *Logger) DisableLevelColor(level int) {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	l.levels[level].color = false
	cf := l.levels[level].colorPrefix
	cf.DisableColor()
	//l.levels[level].SetPrefix(cf.SprintFunc()(prefixes[level]))
}

// EnableLogFormat enables format log of message.
// It will print msg as a log with color or prefix etc.
func (l *Logger) EnableLogFormat() {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	l.EnableColor()
	l.logFormat = true
}

// DisableLogFormat disables format log of message.
// It will print unformatted msg as normally like with any function printf
// without color or prefix etc.
func (l *Logger) DisableLogFormat() {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	l.DisableColor()
	l.logFormat = false
}

// AddHandler add a new handler in the handler list
func (l *Logger) AddHandler(h Handler) {
	mutex.LockOnce()
	defer mutex.UnlockOnce()

	l.handlers = append(l.handlers, h)
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

// Infow logs with info level with structured log format
func (l *Logger) Infow(msg string, v ...interface{}) {
	Log(PRINTW, l, INFO, msg, v...)
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

// Warnw logs with warn level with structured log format
func (l *Logger) Warnw(msg string, v ...interface{}) {
	Log(PRINTW, l, WARN, msg, v...)
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

// Errorw logs with error level with structured log format
func (l *Logger) Errorw(msg string, v ...interface{}) {
	Log(PRINTW, l, ERROR, msg, v...)
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

// Fatalw logs with fatal level with structured log format
func (l *Logger) Fatalw(msg string, v ...interface{}) {
	Log(PRINTW, l, FATAL, msg, v...)
}

//////////// DEFAULT LOGGER ////////////////////////////

// SetVerbosity sets log level. If verbose < NONE, it will be set to NONE.
// If verbose > DEBUG, it will be set to DEBUG
func SetVerbosity(v int) {
	mutex.LockOnce()
	defer mutex.UnlockOnce()

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
	mutex.LockOnce()
	defer mutex.UnlockOnce()

	for i := FATAL; i <= DEBUG; i++ {
		defaultLogger.levels[i].output = w
	}
}

// SetLevelOutput sets output destination for a specific level
func SetLevelOutput(level int, w io.Writer) {
	mutex.LockOnce()
	defer mutex.UnlockOnce()

	defaultLogger.levels[level].output = w
}

// SetFlags sets flags for message log output
func SetFlags(flag int) {
	mutex.LockOnce()
	defer mutex.UnlockOnce()

	defaultLogger.flag = flag
}

// GetFlags gets flags for message log output
func GetFlags() int {
	return defaultLogger.flag
}

// SetTimeFormat sets timestamp with the given format
func SetTimeFormat(format string) {
	mutex.LockOnce()
	defer mutex.UnlockOnce()

	defaultLogger.timeFormat = format
}

// EnableColor enables color for all log levels
func EnableColor() {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	for i := FATAL; i <= DEBUG; i++ {
		defaultLogger.EnableLevelColor(i)
	}
}

// DisableColor disables color for all log levels
func DisableColor() {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	for i := FATAL; i <= DEBUG; i++ {
		defaultLogger.DisableLevelColor(i)
	}
}

// EnableLevelColor enables color for a specific level
func EnableLevelColor(level int) {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	defaultLogger.levels[level].color = true
	cf := defaultLogger.levels[level].colorPrefix
	cf.EnableColor()

}

// DisableLevelColor enables color for a specific level
func DisableLevelColor(level int) {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	defaultLogger.levels[level].color = false
	cf := defaultLogger.levels[level].colorPrefix
	cf.DisableColor()
}

// EnableLogFormat enables format log of message.
// It will print msg as a log with color or prefix etc.
func EnableLogFormat() {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	EnableColor()
	defaultLogger.logFormat = true
}

// DisableLogFormat disables format log of message.
// It will print unformatted msg as normally like with any function printf
// without color or prefix etc.
func DisableLogFormat() {
	if !mutex.IsLocked() {
		mutex.LockOnce()
		defer mutex.UnlockOnce()
	}

	DisableColor()
	defaultLogger.logFormat = false
}

// AddHandler add a new handler in the handler list
func AddHandler(h Handler) {
	mutex.LockOnce()
	defer mutex.UnlockOnce()

	defaultLogger.handlers = append(defaultLogger.handlers, h)
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

// Debugw logs with debug level
func Debugw(msg string, v ...interface{}) {
	Log(PRINTW, defaultLogger, DEBUG, msg, v...)
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

// Infow logs with debug level
func Infow(msg string, v ...interface{}) {
	Log(PRINTW, defaultLogger, INFO, msg, v...)
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

// Warnw logs with debug level
func Warnw(msg string, v ...interface{}) {
	Log(PRINTW, defaultLogger, WARN, msg, v...)
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

// Errorw logs with error level
func Errorw(msg string, v ...interface{}) {
	Log(PRINTW, defaultLogger, ERROR, msg, v...)
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

// Fatalw logs with error level
func Fatalw(msg string, v ...interface{}) {
	Log(PRINTW, defaultLogger, FATAL, msg, v...)
	os.Exit(1)
}

/////////////// INTERNAL FUNCTIONS /////////////////////

// Log wraps print function but using goroutine and waitgroup
// to have a synchronization of logs.
func Log(p int, l *Logger, level int, f string, v ...interface{}) {
	mutex.LockOnce()

	caller := getInfoCaller()
	fields := Fields{}

	fields.Prefix = parsePrefixFields(l, level, caller)
	fields.Log = parseLogFields(p, l, f, v...)

	printMsg(p, l, level, fields)

	defer mutex.UnlockOnce()

	for _, h := range l.handlers {
		h.PrintMsg(p, l, level, fields)
	}
}

func printMsg(p int, l *Logger, level int, fields Fields) {

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

		prefix := formatPrefix(l, cf, fields.Prefix)

		//v := convertMaptoSlice(fields.Log)
		msg := cast.ToString(getField("msg", fields.Log).Value)

		switch p {
		case PRINT:
			if l.logFormat {
				fmt.Fprint(l.levels[level].output, prefix, " ", ct.SprintFunc()(msg))
			} else {
				fmt.Fprint(l.levels[level].output, msg)
			}
		case PRINTF:
			if l.logFormat {
				fmt.Fprintf(l.levels[level].output, "%s %s", prefix, ct.SprintFunc()(msg))
			} else {
				fmt.Fprintf(l.levels[level].output, "%s", msg)
			}
		case PRINTLN:
			if l.logFormat {
				fmt.Fprintln(l.levels[level].output, prefix, ct.SprintFunc()(msg))
			} else {
				fmt.Fprintln(l.levels[level].output, ct.SprintlnFunc()(msg))
			}
		case PRINTW:
			if l.logFormat {
				printw(l, level, prefix, ct, fields.Log)
			} else {
				fmt.Fprintf(l.levels[level].output, "%s\n", ct.SprintfFunc()(msg))
			}
		default:
			if l.logFormat {
				fmt.Fprintln(l.levels[level].output, prefix, ct.SprintFunc()(msg))
			} else {
				fmt.Fprintln(l.levels[level].output, ct.SprintFunc()(msg))
			}
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

func getTimeNow(format string) string {
	return time.Now().Format(format)
}

func quoteString(str string) string {
	s := str
	if strings.Contains(s, " ") {
		s = "\"" + s + "\""
	}

	return s
}

func parsePrefixFields(l *Logger, level int, caller string) []*Field {
	var fields []*Field

	if l.flag&FTIMESTAMP != 0 {
		field := &Field{
			Key:   "ts",
			Value: time.Now().Format(l.timeFormat),
		}
		fields = append(fields, field)
	}

	if l.flag&FCALLER != 0 {
		field := &Field{
			Key:   "caller",
			Value: caller,
		}
		fields = append(fields, field)
	}

	field := &Field{
		Key:   "level",
		Value: prefixes[level],
	}
	fields = append(fields, field)

	return fields
}

func parseLogFields(p int, l *Logger, f string, kv ...interface{}) []*Field {

	var fields []*Field

	switch p {
	case PRINT:
		field := &Field{
			Key:   "msg",
			Value: fmt.Sprint(kv...),
		}
		fields = append(fields, field)
		break
	case PRINTF:
		field := &Field{
			Key:   "msg",
			Value: fmt.Sprintf(f, kv...),
		}
		fields = append(fields, field)
		break
	case PRINTLN:
		field := &Field{
			Key:   "msg",
			Value: fmt.Sprintln(kv...),
		}
		fields = append(fields, field)

	case PRINTW:
		field := &Field{
			Key:   "msg",
			Value: f,
		}
		fields = append(fields, field)

		if len(kv)%2 != 0 {
			kv = append(kv, "missing")
		}

		// if no key/value fields, return line after print message

		for i := 0; i < len(kv); i += 2 {
			// cast 1st elem = key to string
			k := cast.ToString(kv[i])
			if k == "" {
				k = "missing"
			}

			// cast 2nd elem = value
			var val string
			v := kv[i+1]
			kind := reflect.ValueOf(v).Kind()
			if kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map || kind == reflect.Struct || kind == reflect.Ptr {
				val = fmt.Sprintf("%+v", v)
			} else {
				s := cast.ToString(v)
				val = quoteString(s)
			}
			field := &Field{
				Key:   k,
				Value: val,
			}
			fields = append(fields, field)
		}
		break
	}

	return fields
}

func getField(key string, fields []*Field) *Field {
	for _, f := range fields {
		if f.Key == key {
			return f
		}
	}

	return nil
}

func printw(l *Logger, level int, prefix string, ct *color.Color, fields []*Field) {
	var pairs []interface{}
	var format string
	var message string

	output := l.levels[level].output

	msg := cast.ToString(getField("msg", fields).Value)
	if l.flag&FFULLSTRUCTUREDLOG != 0 {
		message = ct.SprintFunc()("msg=") + quoteString(msg)
		format += "%s %s "
	} else {
		message = ct.SprintFunc()(msg)
		format += "%s %-60s "
	}

	pairs = append(pairs, prefix, message)

	// Delete msg once used
	dupFields := append([]*Field{}, fields[1:]...)

	// if no key/value fields, return line after print message
	if len(dupFields) <= 0 {
		format += "\n"
	} else {

		// Tip to keep order while parsing fields
		for i, field := range dupFields {
			// key
			k := field.Key

			// value
			v := field.Value
			pair := ""
			pair = fmt.Sprintf("%s=%s", ct.SprintFunc()(k), v)

			//pair = fmt.Sprintf("%s=%+v", k, v)
			pairs = append(pairs, pair)
			if i < len(dupFields)-1 {
				format += "%s "
			} else {
				format += "%s\n"
			}
		}
	}

	fmt.Fprintf(output, format, pairs...)
}

func formatPrefix(l *Logger, cf *color.Color, fields []*Field) string {
	var format string
	var values []interface{}

	if l.flag&FTIMESTAMP != 0 {
		ts := getField("ts", fields).Value
		format += "%s "
		if l.flag&FFULLSTRUCTUREDLOG != 0 {
			values = append(values, "ts="+ts)
		} else {
			values = append(values, ts)
		}
	}

	if l.flag&FCALLER != 0 {
		caller := getField("caller", fields).Value
		format += "%s "
		if l.flag&FFULLSTRUCTUREDLOG != 0 {
			values = append(values, "caller="+caller)
		} else {
			values = append(values, caller)
		}
	}

	level := getField("level", fields).Value
	if l.flag&FFULLSTRUCTUREDLOG != 0 {
		format += "%s"
		values = append(values, "level="+cf.SprintFunc()(level))
	} else {
		format += "%-17s"
		values = append(values, cf.SprintFunc()(level+":"))
	}

	return fmt.Sprintf(format, values...)
}
