package logger

import (
	"fmt"
	"strings"
	"time"
)

func Info(msg string) {
	fmt.Printf("%s :: [INFO] %s\n", time.Now().Format(time.RFC3339), strings.TrimSpace(msg))
}

func Infof(format string, a ...any) {
	fmt.Printf("%s :: [INFO] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...))
}

func Error(msg string) {
	fmt.Printf("%s :: [ERROR] %s\n", time.Now().Format(time.RFC3339), strings.TrimSpace(msg))
}

func Errorf(format string, a ...any) {
	fmt.Printf("%s :: [ERROR] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, a...))
}
