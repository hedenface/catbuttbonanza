package log

import (
    "testing"
    "log"
    "bytes"
    "time"
    "fmt"
    "io"
)

const (
    expectedData = "test"
)

type funcPrintln func(...any)
type funcPrintf func(string, ...any)
type funcPrint func()

func getExpectedOutput(l int, app string, data string) string {
    return fmt.Sprintf("%s [%s] [%s] ([test]) [%s]\n", time.Now().Format("2006/01/02 15:04:05"), app, GetLevelString(l), data)
}

func testTraceNoOutput(buf *bytes.Buffer, t *testing.T) {
    PushStack("test")
    defer PopStack()

    buf.Reset()

    Trace("%s", expectedData)
 
    _, err := buf.ReadString('\n')

    // in this case, err should not be nil, since we're expecting EOF
    if err != io.EOF {
        t.Fatalf("Expected ReadString('\\n') on buf to have err = EOF, but '%v' was the actual error.", err)
    }
}

func testPrint(l int, f funcPrint, t *testing.T) {
    PushStack("test")
    defer PopStack()

    app := "a"

    Setup(app, l)

    var buf bytes.Buffer

    logger.SetOutput(&buf)

    // JUST IN CASE the second is different, so we check against both "expected outputs"
    expectedOutputWithBeforeTime := getExpectedOutput(l, app, expectedData)
    f()
    expectedOutputWithAfterTime := getExpectedOutput(l, app, expectedData)
 
    loggedOutput, err := buf.ReadString('\n')

    if err != nil {
        t.Fatalf("Unable to call ReadString('\\n') on buf, error: %v.", err)
    }

    if loggedOutput != expectedOutputWithBeforeTime &&
       loggedOutput != expectedOutputWithAfterTime {
        t.Fatalf("Expected output of either \n'%s' or \n'%s', \n'%s' was logged instead.", expectedOutputWithBeforeTime, expectedOutputWithAfterTime, loggedOutput)
    }

    // now we check if it's a greater level than Trace, to see if the level below actually shows up in the buffer
    if l == TraceLevel {
        return
    }

    testTraceNoOutput(&buf, t)
}

func TestSetup(t *testing.T) {
    expectedPrefix := "[name] "
    expectedFlags := log.LstdFlags|log.Lmsgprefix

    Setup("name", InfoLevel)

    if Level != InfoLevel {
        t.Fatalf("Level Set Improperly. Needed InfoLevel=%d, %d was set instead.", InfoLevel, Level)
    }

    if logger.Prefix() != expectedPrefix {
        t.Fatalf("Logger prefix set improperly. Needed '%s', '%s' was set instead.", expectedPrefix, logger.Prefix())
    }

    if logger.Flags() != expectedFlags {
        t.Fatalf("Logger flags set improperly. Needed %d, %d was set instead.", expectedFlags, logger.Flags())
    }
}

func testPrintln(level int, f funcPrintln, t *testing.T) {
    testPrint(level, func() { f(expectedData) }, t)
}

func testPrintf(level int, f funcPrintf, t *testing.T) {
    testPrint(level, func() { f("%s", expectedData) }, t)
}

func TestTraceln(t *testing.T) {
    testPrintln(TraceLevel, Traceln, t)
}

func TestTrace(t *testing.T) {
    testPrintf(TraceLevel, Trace, t)
}

func TestDebugln(t *testing.T) {
    testPrintln(DebugLevel, Debugln, t)
}

func TestDebug(t *testing.T) {
    testPrintf(DebugLevel, Debug, t)
}

func TestInfoln(t *testing.T) {
    testPrintln(InfoLevel, Infoln, t)
}

func TestInfo(t *testing.T) {
    testPrintf(InfoLevel, Info, t)
}

func TestWarnln(t *testing.T) {
    testPrintln(WarnLevel, Warnln, t)
}

func TestWarn(t *testing.T) {
    testPrintf(WarnLevel, Warn, t)
}

func TestErrorln(t *testing.T) {
    testPrintln(ErrorLevel, Errorln, t)
}

func TestError(t *testing.T) {
    testPrintf(ErrorLevel, Error, t)
}

func TestFatalln(t *testing.T) {
    Exiter = func(int) { }
    testPrintln(FatalLevel, Fatalln, t)
}

func TestFatal(t *testing.T) {
    Exiter = func(int) {  }
    testPrintf(FatalLevel, Fatal, t)
}

func testGetLevelString(l int, s string, t *testing.T) {
    if GetLevelString(l) != s {
        t.Fatalf("GetLevelString(%d) was expected to be '%s', '%s' was returned instead.", l, s, GetLevelString(l))
    }
}

func TestGetLevelString(t *testing.T) {
    testGetLevelString(TraceLevel, "TRACE", t)
    testGetLevelString(DebugLevel, "DEBUG", t)
    testGetLevelString(InfoLevel, "INFO", t)
    testGetLevelString(WarnLevel, "WARN", t)
    testGetLevelString(ErrorLevel, "ERROR", t)
    testGetLevelString(FatalLevel, "FATAL", t)
    testGetLevelString(42, "FATAL", t)
}

func TestPushStack(t *testing.T) {
    PushStack("test")
    if len(funcStack) != 1 {
        t.Fatalf("PushStack() failed. Expected funcStack length of 1, %d was returned instead.", len(funcStack))
    }
    PopStack()
}

func TestPopStack(t *testing.T) {
    PushStack("test")
    PopStack()
    if len(funcStack) != 0 {
        t.Fatalf("PopStack() failed. Expected funcStack length of 0, %d was returned instead.", len(funcStack))
    }
}
