package golog_test

import (
	"bytes"
	"io"
	"os"
	//"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/uthng/golog"
	utils "github.com/uthng/goutils"
)

func TestSimpleLog(t *testing.T) {
	testCases := []struct {
		name    string
		verbose int
		output  []string
	}{
		{
			"Debug",
			golog.DEBUG,
			[]string{
				`DEBUG:[ ]+This is debug log$`,
				`INFO:[ ]+This is info log$`,
				`WARN:[ ]+This is warn log$`,
				`ERROR:[ ]+This is error log$`,
			},
		},
		{
			"Info",
			golog.INFO,
			[]string{
				`INFO:[ ]+This is info log$`,
				`WARN:[ ]+This is warn log$`,
				`ERROR:[ ]+This is error log$`,
			},
		},
		{
			"Warn",
			golog.WARN,
			[]string{
				`WARN:[ ]+This is warn log$`,
				`ERROR:[ ]+This is error log$`,
			},
		},
		{
			"Error",
			golog.ERROR,
			[]string{
				`ERROR:[ ]+This is error log$`,
			},
		},
		{
			"None",
			golog.NONE,
			[]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			multi := io.MultiWriter(&buf, os.Stdout)
			logger := golog.NewLogger()
			logger.SetVerbosity(tc.verbose)
			logger.SetOutput(multi)

			logger.Debugln("This is debug log")
			logger.Infoln("This is info log")
			logger.Warnln("This is warn log")
			logger.Errorln("This is error log")

			if logger.GetVerbosity() != golog.NONE {
				str := utils.StripAnsi(buf.String())
				arr := strings.Split(str, "\n\n")
				for idx, w := range tc.output {
					matched, _ := regexp.MatchString(w, arr[idx])
					if !matched {
						t.Errorf("\nwant:\n%s\nhave:\n%s", w, arr[idx])
					}
				}
			} else {
				if len(buf.String()) != 0 {
					t.Errorf("\nwant:\n%s\nhave:\n%s", strings.Join(tc.output, ""), buf.String())
				}
			}
		})
	}

}

func TestFormattedLog(t *testing.T) {
	testCases := []struct {
		name    string
		verbose int
		output  []string
	}{
		{
			"Debug",
			golog.DEBUG,
			[]string{
				`DEBUG:[ ]+This is debug log$`,
				`INFO:[ ]+This is info log$`,
				`WARN:[ ]+This is warn log$`,
				`ERROR:[ ]+This is error log$`,
			},
		},
		{
			"Info",
			golog.INFO,
			[]string{
				`INFO:[ ]+This is info log$`,
				`WARN:[ ]+This is warn log$`,
				`ERROR:[ ]+This is error log$`,
			},
		},
		{
			"Warn",
			golog.WARN,
			[]string{
				`WARN:[ ]+This is warn log$`,
				`ERROR:[ ]+This is error log$`,
			},
		},
		{
			"Error",
			golog.ERROR,
			[]string{
				`ERROR:[ ]+This is error log$`,
			},
		},
		{
			"None",
			golog.NONE,
			[]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			multi := io.MultiWriter(&buf, os.Stdout)
			logger := golog.NewLogger()
			logger.SetVerbosity(tc.verbose)
			logger.SetOutput(multi)

			logger.Debugf("This is %s log\n", "debug")
			logger.Infof("This is %s log\n", "info")
			logger.Warnf("This is %s log\n", "warn")
			logger.Errorf("This is %s log\n", "error")

			if logger.GetVerbosity() != golog.NONE {
				arr := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n\n"), []byte("\n"))
				for idx, w := range tc.output {
					msg := utils.StripAnsi(string(arr[idx]))
					matched, _ := regexp.MatchString(w, msg)
					if !matched {
						t.Errorf("\nwant:\n%s\nhave:\n%s", w, msg)
					}
				}
			} else {
				if len(buf.String()) != 0 {
					t.Errorf("\nwant:\n%s\nhave:\n%s", strings.Join(tc.output, ""), buf.String())
				}
			}
		})
	}

}

func TestLogColor(t *testing.T) {
	testCases := []struct {
		name   string
		level  int
		color  bool
		output string
	}{
		{
			"Debug",
			golog.DEBUG,
			false,
			`DEBUG:[ ]+This is debug log`,
		},
		{
			"Info",
			golog.INFO,
			true,
			`INFO:[ ]+This is info log`,
		},
		{
			"Warn",
			golog.WARN,
			false,
			`WARN:[ ]+This is warn log`,
		},
		{
			"Error",
			golog.ERROR,
			true,
			`ERROR:[ ]+This is error log`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			multi := io.MultiWriter(&buf, os.Stdout)
			logger := golog.NewLogger()
			logger.SetVerbosity(5)
			logger.SetOutput(multi)

			if tc.color {
				logger.EnableLevelColor(tc.level)
			} else {
				logger.DisableLevelColor(tc.level)
			}

			if tc.level == golog.DEBUG {
				logger.Debugf("This is %s log\n", "debug")
			} else if tc.level == golog.INFO {
				logger.Infof("This is %s log\n", "info")
			} else if tc.level == golog.WARN {
				logger.Warnf("This is %s log\n", "warn")
			} else {
				logger.Errorf("This is %s log\n", "error")
			}

			b := bytes.TrimRight(buf.Bytes(), "\n\n")
			msg := string(b)
			if tc.color {
				msg = utils.StripAnsi(msg)
			}
			matched, _ := regexp.MatchString(tc.output, msg)
			if !matched {
				t.Errorf("\nwant:\n%s\nhave:\n%s", tc.output, msg)
			}
		})
	}

}

func TestLogColorAll(t *testing.T) {
	testCases := []struct {
		name   string
		color  bool
		output []string
	}{
		{
			"EnableAll",
			true,
			[]string{
				`DEBUG:[ ]+This is debug log`,
				`INFO:[ ]+This is info log`,
				`WARN:[ ]+This is warn log`,
				`ERROR:[ ]+This is error log`,
			},
		},
		{
			"DisableAll",
			false,
			[]string{
				`DEBUG:[ ]+This is debug log`,
				`INFO:[ ]+This is info log`,
				`WARN:[ ]+This is warn log`,
				`ERROR:[ ]+This is error log`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			multi := io.MultiWriter(&buf, os.Stdout)
			logger := golog.NewLogger()
			logger.SetVerbosity(5)
			logger.SetOutput(multi)

			if tc.color {
				logger.EnableColor()
			} else {
				logger.DisableColor()
			}

			logger.Debugf("This is %s log\n", "debug")
			logger.Infof("This is %s log\n", "info")
			logger.Warnf("This is %s log\n", "warn")
			logger.Errorf("This is %s log\n", "error")

			arr := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n\n"), []byte("\n"))
			for idx, w := range tc.output {
				msg := string(arr[idx])
				if tc.color {
					msg = utils.StripAnsi(msg)
				}
				matched, _ := regexp.MatchString(w, msg)
				if !matched {
					t.Errorf("\nwant:\n%s\nhave:\n%s", w, msg)
				}
			}
		})
	}
}

func TestLogDefault(t *testing.T) {
	testCases := []struct {
		name   string
		output string
	}{
		{
			"Error",
			`ERROR:[ ]+This is error log`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			multi := io.MultiWriter(&buf, os.Stdout)
			logger := golog.NewLogger()
			logger.SetVerbosity(5)
			logger.SetLevelOutput(golog.ERROR, multi)

			logger.Debugf("This is %s log\n", "debug")
			logger.Infof("This is %s log\n", "info")
			logger.Warnf("This is %s log\n", "warn")
			logger.Errorf("This is %s log\n", "error")

			b := bytes.TrimRight(buf.Bytes(), "\n\n")
			msg := string(b)
			msg = utils.StripAnsi(msg)

			matched, _ := regexp.MatchString(tc.output, msg)
			if !matched {
				t.Errorf("\nwant:\n%s\nhave:\n%s", tc.output, msg)
			}
		})
	}

}

func TestLogCaller(t *testing.T) {
	var buf bytes.Buffer

	output := []string{
		`(.*) log_test.go:365:TestLogCaller DEBUG:[ ]+This is debug log$`,
		`(.*) log_test.go:366:TestLogCaller INFO:[ ]+This is info log$`,
		`(.*) log_test.go:367:TestLogCaller WARN:[ ]+This is warn log$`,
		`(.*) log_test.go:368:TestLogCaller ERROR:[ ]+This is error log$`,
	}

	multi := io.MultiWriter(&buf, os.Stdout)
	logger := golog.NewLogger()
	logger.SetOutput(multi)
	logger.SetVerbosity(5)
	logger.SetFlags(golog.FTIMESTAMP | golog.FCALLER)

	logger.Debugf("This is %s log\n", "debug")
	logger.Infof("This is %s log\n", "info")
	logger.Warnf("This is %s log\n", "warn")
	logger.Errorf("This is %s log\n", "error")

	arr := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n\n"), []byte("\n"))
	for idx, w := range output {
		msg := string(arr[idx])
		msg = utils.StripAnsi(msg)

		matched, _ := regexp.MatchString(w, msg)
		if !matched {
			t.Errorf("\nwant:\n%s\nhave:\n%s", w, msg)
		}
	}

}

func TestLogWith(t *testing.T) {
	testCases := []struct {
		name   string
		flag   int
		output []string
	}{
		{
			"NoFullStructured",
			golog.FTIMESTAMP | golog.FCALLER,
			[]string{
				`(.*) log_test.go:422:func1 DEBUG:[ ]+This is debug log[ ]+level="debug level" value=15.5`,
				`(.*) log_test.go:423:func1 INFO:[ ]+This is info log[ ]+level="info level" value=15.5`,
				`(.*) log_test.go:424:func1 WARN:[ ]+This is warn log[ ]+level="warn level" value=15.5`,
				`(.*) log_test.go:425:func1 ERROR:[ ]+This is error log[ ]+level="error level" value=15.5`,
			},
		},
		{
			"FullStructured",
			golog.FTIMESTAMP | golog.FCALLER | golog.FFULLSTRUCTUREDLOG,
			[]string{
				`ts=(.*) caller=log_test.go:422:func1 level=DEBUG msg="This is debug log" level="debug level" value=15.5`,
				`ts=(.*) caller=log_test.go:423:func1 level=INFO msg="This is info log" level="info level" value=15.5`,
				`ts=(.*) caller=log_test.go:424:func1 level=WARN msg="This is warn log" level="warn level" value=15.5`,
				`ts=(.*) caller=log_test.go:425:func1 level=ERROR msg="This is error log" level="error level" value=15.5`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			multi := io.MultiWriter(&buf, os.Stdout)
			logger := golog.NewLogger()
			logger.SetOutput(multi)
			logger.SetVerbosity(5)
			logger.SetFlags(tc.flag)
			logger.SetTimeFormat("2006-01-02T15:04:05.000000")

			logger.Debugw("This is debug log", "level", "debug level", "value", 15.5)
			logger.Infow("This is info log", "level", "info level", "value", 15.5)
			logger.Warnw("This is warn log", "level", "warn level", "value", 15.5)
			logger.Errorw("This is error log", "level", "error level", "value", 15.5)

			arr := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n\n"), []byte("\n"))
			for idx, w := range tc.output {
				msg := string(arr[idx])
				msg = utils.StripAnsi(msg)

				matched, _ := regexp.MatchString(w, msg)
				if !matched {
					t.Errorf("\nwant:\n%s\nhave:\n%s", w, msg)
				}
			}
		})
	}
}

func TestLogFormat(t *testing.T) {
	var buf bytes.Buffer

	output := []string{
		`This is debug log`,
		`This is info log`,
		`This is warn log`,
		`This is error log`,
	}

	multi := io.MultiWriter(&buf, os.Stdout)
	logger := golog.NewLogger()
	logger.SetOutput(multi)
	logger.SetVerbosity(5)
	logger.DisableLogFormat()

	// Test if any kind of print will produce the same unformatted message
	logger.Debugw("This is debug log")
	logger.Infow("This is info log")
	logger.Warnf("This is warn log\n")
	logger.Error("This is error log\n")

	arr := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n\n"), []byte("\n"))
	for idx, w := range output {
		msg := string(arr[idx])

		matched, _ := regexp.MatchString(w, msg)
		if !matched {
			t.Errorf("\nwant:\n%s\nhave:\n%s", w, msg)
		}
	}

}
