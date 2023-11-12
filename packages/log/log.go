package helper

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"strings"
)

const (
	Trace int = 0
	Debug int = 1
	Info  int = 2
	Warn  int = 3
	Error int = 4
)

var (
	Level int = Info
)

func SetupLogger(pkg string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", pkg), log.LstdFlags|log.Lmsgprefix)
}

func SetLoggerLevel(level string) {
	Level = getNumericLevel(level)
}

func (l *log.Logger) Log(level string, message string) {
	if getNumericLevel(level) >= Level {
		l.Printf("[%s] %s", strings.ToUpper(level), message)
	}
}

//func (l *log.Logger) Trace()

func getNumericLevel(level string) int {
	switch strings.ToUpper(level) {
	case "TRACE":
		return Trace
	case "DEBUG":
		return Debug
	case "INFO":
		return Info
	case "WARN":
		return Warn
	case "ERROR":
		return Error
	}

	return -1
}
