package log

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
)

var (
	white  = color.New(color.FgWhite).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
)

func Info(format string, args ...interface{}) {

	t := time.Now().Format("2006-01-02 15:04:05")
	pc, file, line, _ := runtime.Caller(1)

	prefix := fmt.Sprintf("%s %s:%s:%d:", t, file, runtime.FuncForPC(pc).Name(), line)
	content := fmt.Sprintf(format, args...)

	fmt.Println(blue(prefix), white(content))
}

func Debug(format string, args ...interface{}) {

	t := time.Now().Format("2006-01-02 15:04:05")
	pc, file, line, _ := runtime.Caller(1)

	prefix := fmt.Sprintf("%s %s:%s:%d:", t, file, runtime.FuncForPC(pc).Name(), line)
	content := fmt.Sprintf(format, args...)

	fmt.Println(green(prefix), white(content))
}

func Warn(format string, args ...interface{}) {

	t := time.Now().Format("2006-01-02 15:04:05")
	pc, file, line, _ := runtime.Caller(1)

	prefix := fmt.Sprintf("%s %s:%s:%d:", t, file, runtime.FuncForPC(pc).Name(), line)
	content := fmt.Sprintf(format, args...)

	fmt.Println(yellow(prefix), white(content))
}

func Error(format string, args ...interface{}) {

	t := time.Now().Format("2006-01-02 15:04:05")
	pc, file, line, _ := runtime.Caller(1)

	prefix := fmt.Sprintf("%s %s:%s:%d:", t, file, runtime.FuncForPC(pc).Name(), line)
	content := fmt.Sprintf(format, args...)

	fmt.Println(red(prefix), white(content))
}

func Fatal(format string, args ...interface{}) {

	t := time.Now().Format("2006-01-02 15:04:05")
	pc, file, line, _ := runtime.Caller(1)

	prefix := fmt.Sprintf("%s %s:%s:%d:", t, file, runtime.FuncForPC(pc).Name(), line)
	content := fmt.Sprintf(format, args...)

	fmt.Println(red(prefix), white(content))
	
	// will only exit in this level because it is fatal
	os.Exit(1)
}

