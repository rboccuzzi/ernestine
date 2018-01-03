package log

import (
  "fmt"
  "os"
)

const (
  fatal = -1
  error = iota
  info  = iota
  debug = iota
)

var verbosityLevel int

// SetVerbosityLevel sets the level that controls what messages get output to log
func SetVerbosityLevel(level int) {
  verbosityLevel = level
}

// Log takes level of message and outputs if VerbosityLevel is set above given level
func Log(level int, message ...interface{}) {
  if level <= verbosityLevel {
    fmt.Println(message...)
  }
}

func Fatal(message ...interface{}) {
  Log(fatal, append([]interface{}{"FATAL ERROR:"}, message...)...)
  os.Exit(-1)
}

func Error(message ...interface{}) {
  Log(error, message...)
}

func Debug(message ...interface{}) {
  Log(debug, message...)
}

func Info(message ...interface{}) {
  Log(info, message...)
}
