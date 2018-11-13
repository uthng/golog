package golog_test

import (
	"bytes"
	"io"
	"os"
	//"reflect"
	//"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/uthng/common/utils"
	"github.com/uthng/golog"
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
				`DEBUG: (.*) This is debug log$`,
				`INFO: (.*) This is info log$`,
				`WARN: (.*) This is warn log$`,
				`ERROR: (.*) This is error log$`,
			},
		},
		{
			"Info",
			golog.INFO,
			[]string{
				`INFO: (.*) This is info log$`,
				`WARN: (.*) This is warn log$`,
				`ERROR: (.*) This is error log$`,
			},
		},
		{
			"Warn",
			golog.WARN,
			[]string{
				`WARN: (.*) This is warn log$`,
				`ERROR: (.*) This is error log$`,
			},
		},
		{
			"Error",
			golog.ERROR,
			[]string{
				`ERROR: (.*) This is error log$`,
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
			logger := golog.NewLogger(multi)
			logger.SetVerbosity(tc.verbose)

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
				`DEBUG: (.*) This is debug log$`,
				`INFO: (.*) This is info log$`,
				`WARN: (.*) This is warn log$`,
				`ERROR: (.*) This is error log$`,
			},
		},
		{
			"Info",
			golog.INFO,
			[]string{
				`INFO: (.*) This is info log$`,
				`WARN: (.*) This is warn log$`,
				`ERROR: (.*) This is error log$`,
			},
		},
		{
			"Warn",
			golog.WARN,
			[]string{
				`WARN: (.*) This is warn log$`,
				`ERROR: (.*) This is error log$`,
			},
		},
		{
			"Error",
			golog.ERROR,
			[]string{
				`ERROR: (.*) This is error log$`,
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
			logger := golog.NewLogger(multi)
			logger.SetVerbosity(tc.verbose)

			logger.Debugf("This is %s log", "debug")
			logger.Infof("This is %s log", "info")
			logger.Warnf("This is %s log", "warn")
			logger.Errorf("This is %s log", "error")

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
