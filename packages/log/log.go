package log

import (
	"fmt"
	"log"
	"os"
	"strings"
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
type FunctionStack []string

var (
	Level int
	logger *log.Logger
	Exiter ExiterFunction
	DefaultExitor ExiterFunction
	funcStack FunctionStack
)

func (fs FunctionStack) String() string {
	s := ""
	for _, f := range fs {
		s = fmt.Sprintf("%s[%s] ", s, f)
	}
	if s != "" {
		s = fmt.Sprintf("(%s)", strings.TrimSuffix(s, " "))
	}
	return s
}

func Setup(app string, level int) {
	Level = level
	logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", app), log.LstdFlags|log.Lmsgprefix)

	// we set this here, so that we can override it in the test suite
	DefaultExitor = os.Exit

	cbb, err := os.ReadFile("resources/buttbonanza")
	if err != nil {
		fmt.Println("Could not read buttbonanza file...")
		fmt.Println("", "", "Cat Butt Bonanza")
	} else {
		fmt.Println(string(cbb))
	}
	// we want this centered. the text is centered at character 15
	if len(app) > 30 {
		fmt.Println(app)
	} else {
		startPoint := 15 + (len(app) / 2)
		fmtString := fmt.Sprintf("%%%ds\n", startPoint)
		fmt.Printf(fmtString, app)
	}
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

func PushStack(f string) {
	funcStack = append(funcStack, f)
}

func PopStack() {
	n := len(funcStack)-1
	funcStack[n] = ""
	funcStack = funcStack[:n]
}

func Println(l int, a ...any) {
	if l >= Level {
		logger.Println(fmt.Sprintf("[%s]", GetLevelString(l)), funcStack.String(), a)
	}
}

func Printf(l int, format string, a ...any) {
	if l >= Level {
		newFormat := fmt.Sprintf("[%s] %s %s", GetLevelString(l), funcStack.String(), format)
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