package log

import (
	"fmt"
	"log"
	"os"
)

const (
	TraceLevel = 0
	DebugLevel = 1
	InfoLevel  = 2
	WarnLevel  = 3
	ErrorLevel = 4
	FatalLevel = 5
)

type ExiterFunction func(int)

var (
	Level int
	logger *log.Logger
	Exiter ExiterFunction
	DefaultExitor ExiterFunction
)

func Setup(app string, level int) {
	Level = level
	logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", app), log.LstdFlags|log.Lmsgprefix)

	// we set this here, so that we can override it in the test suite
	DefaultExitor = os.Exit
}

func GetLevelString(l int) string {
	switch l {
    case TraceLevel:
        return "TRACE"
    case DebugLevel:
        return "DEBUG"
    case InfoLevel:
        return "INFO"
    case WarnLevel:
        return "WARN"
    case ErrorLevel:
        return "ERROR"
    }

    return "FATAL"
}

func Println(l int, a ...any) {
	if l >= Level {
		logger.Println(fmt.Sprintf("[%s]", GetLevelString(l)), a)
	}
}

func Printf(l int, format string, a ...any) {
	if l >= Level {
		newFormat := fmt.Sprintf("[%s] %s", GetLevelString(l), format)
		logger.Printf(newFormat, a)
	}
}

func Traceln(a ...any) {
	Println(TraceLevel, a...)
}

func Trace(format string, a ...any) {
	Printf(TraceLevel, format, a...)
}

func Debugln(a ...any) {
	Println(DebugLevel, a...)
}

func Debug(format string, a ...any) {
	Printf(DebugLevel, format, a...)
}

func Infoln(a ...any) {
	Println(InfoLevel, a...)
}

func Info(format string, a ...any) {
	Printf(InfoLevel, format, a...)
}

func Warnln(a ...any) {
	Println(WarnLevel, a...)
}

func Warn(format string, a ...any) {
	Printf(WarnLevel, format, a...)
}

func Errorln(a ...any) {
	Println(ErrorLevel, a...)
}

func Error(format string, a ...any) {
	Printf(ErrorLevel, format, a...)
}

func Fatalln(a ...any) {
	Println(FatalLevel, a...)
	if Exiter == nil {
		Exiter = DefaultExitor
	}
	Exiter(1)
}

func Fatal(format string, a ...any) {
	Printf(FatalLevel, format, a...)
	if Exiter == nil {
		Exiter = DefaultExitor
	}
	Exiter(1)
}