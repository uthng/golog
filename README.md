# golog
Simple logging library using golang log package. It uses io.Writer as output destination so it can log to stdout, stderr
or into a file or multiple destinations at the same time using io.MultiWriter.

## Documentation
See the [Godoc](https://godoc.org/github.com/uthng/golog)

## Usage

#### Set up a standard logger to stdout:

```
package main

import (
  "os"
  "github.com/uthng/common/golog"
)

function main() {
  logger := golog.NewLogger()
  logger.SetVerbosity(golog.DEBUG)

  logger.Debugln("This is debug log")
  logger.Infoln("This is info log")
  logger.Warnln("This is warn log")
  logger.Errorln("This is error log")
}
```

And when executed, the program will show the following output to the standard output:

```
DEBUG: 2018/11/12 02:08:57 log.go:161: This is debug log
INFO: 2018/11/12 02:08:57 log.go:161: This is info log
WARN: 2018/11/12 02:08:57 log.go:161: This is warn log
ERROR: 2018/11/12 02:08:57 log.go:161: This is error log
```

#### Set up multiple destinations
```
package main

import (
  "io"
  "os"
  "github.com/uthng/common/golog"
)

function main() {
  file, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  if err != nil {
	  os.Exit(-1)
  }

  multi := io.MultiWriter(file, os.Stdout)

  logger := golog.NewLogger()
  logger.SetVerbosity(golog.INFO)
  logger.SetOutput(multi)
  
  logger.Debugln("This is debug log")
  logger.Infoln("This is info log")
  logger.Warnln("This is warn log")
  logger.Errorln("This is error log")
}
```

When executed, the program will log to stdout and into file.txt the following output:

```
INFO: 2018/11/12 02:19:13 log.go:161: This is info log
WARN: 2018/11/12 02:19:13 log.go:161: This is warn log
ERROR: 2018/11/12 02:19:13 log.go:161: This is error log
```
