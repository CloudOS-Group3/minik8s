package log

import (
	"fmt"
	"os"
	"path/filepath"
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

// set this value to output all log.Debug()
const isDebug = true

// usage; you can use these log functions the same way you use fmt.Printf
// the differece is that all the functions below will add a new line at the end

func Info(format string, args ...interface{}) {

	t := time.Now().Format("2006-01-02 15:04:05")
	pc, targetPath, line, _ := runtime.Caller(1)

	absPath, err := os.Getwd()
	if err != nil {
		return
	}

	relPath, err := filepath.Rel(absPath, targetPath)
	if err != nil {
		return
	}

	prefix := fmt.Sprintf("%s %s:%s:%d:", t, relPath, runtime.FuncForPC(pc).Name(), line)
	content := fmt.Sprintf(format, args...)

	fmt.Println(blue(prefix), white(content))
}

func Debug(format string, args ...interface{}) {

	if !isDebug {
		return
	}

	t := time.Now().Format("2006-01-02 15:04:05")
	pc, targetPath, line, _ := runtime.Caller(1)

	absPath, err := os.Getwd()
	if err != nil {
		return
	}

	relPath, err := filepath.Rel(absPath, targetPath)
	if err != nil {
		return
	}

	prefix := fmt.Sprintf("%s %s:%s:%d:", t, relPath, runtime.FuncForPC(pc).Name(), line)
	content := fmt.Sprintf(format, args...)

	fmt.Println(green(prefix), white(content))
}

func Warn(format string, args ...interface{}) {

	t := time.Now().Format("2006-01-02 15:04:05")
	pc, targetPath, line, _ := runtime.Caller(1)

	absPath, err := os.Getwd()
	if err != nil {
		return
	}

	relPath, err := filepath.Rel(absPath, targetPath)
	if err != nil {
		return
	}

	prefix := fmt.Sprintf("%s %s:%s:%d:", t, relPath, runtime.FuncForPC(pc).Name(), line)
	content := fmt.Sprintf(format, args...)

	fmt.Println(yellow(prefix), white(content))
}

func Error(format string, args ...interface{}) {

	t := time.Now().Format("2006-01-02 15:04:05")
	pc, targetPath, line, _ := runtime.Caller(1)

	absPath, err := os.Getwd()
	if err != nil {
		return
	}

	relPath, err := filepath.Rel(absPath, targetPath)
	if err != nil {
		return
	}

	prefix := fmt.Sprintf("%s %s:%s:%d:", t, relPath, runtime.FuncForPC(pc).Name(), line)
	content := fmt.Sprintf(format, args...)

	fmt.Println(red(prefix), white(content))
}

func Fatal(format string, args ...interface{}) {

	t := time.Now().Format("2006-01-02 15:04:05")
	pc, targetPath, line, _ := runtime.Caller(1)

	absPath, err := os.Getwd()
	if err != nil {
		return
	}

	relPath, err := filepath.Rel(absPath, targetPath)
	if err != nil {
		return
	}

	prefix := fmt.Sprintf("%s %s:%s:%d:", t, relPath, runtime.FuncForPC(pc).Name(), line)
	content := fmt.Sprintf(format, args...)

	fmt.Println(red(prefix), white(content))

	// will only exit in this level because it is fatal
	os.Exit(1)
}
