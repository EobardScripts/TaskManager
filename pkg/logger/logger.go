package logger

import (
	"fmt"
	"runtime"
	"time"
)

var Console = true

func getFuncName() string {
	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return details.Name()
	}
	return "undefined"
}

func New() error {
	return nil
}

func Info(format string, a ...any) {
	date := time.Now()
	dateFormat := fmt.Sprintf("[%s]", date.Format("02.01.2006 15.01"))
	funcName := fmt.Sprintf("%s ->", getFuncName())
	str := fmt.Sprintf("%s %s info: %s", dateFormat, funcName, format)
	if Console {
		fmt.Printf(str, a...)
		fmt.Println()
	}
}

func Error(format string, a ...any) {
	date := time.Now()
	dateFormat := fmt.Sprintf("[%s]", date.Format("02.01.2006 15.01"))
	funcName := fmt.Sprintf("%s ->", getFuncName())
	str := fmt.Sprintf("%s %s error: %s", dateFormat, funcName, format)
	if Console {
		fmt.Printf(str, a...)
		fmt.Println()
	}
}

func Warn(format string, a ...any) {
	date := time.Now()
	dateFormat := fmt.Sprintf("[%s]", date.Format("02.01.2006 15.01"))
	funcName := fmt.Sprintf("%s ->", getFuncName())
	str := fmt.Sprintf("%s %s warn: %s", dateFormat, funcName, format)
	if Console {
		fmt.Printf(str, a...)
		fmt.Println()
	}
}
