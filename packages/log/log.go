package log

import (
	"fmt"
	"log"
	"os"
)

func Setup(app string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", app), log.LstdFlags|log.Lmsgprefix)
}
