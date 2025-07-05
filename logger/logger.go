package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	loggerVar *log.Logger
)

type Color int8

const (
	WHITE Color = iota
	YELLOW
	RED
)

func init() {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	loggerVar = log.New(multiWriter, "", 0)
}

// https://www.dolthub.com/blog/2024-02-23-colors-in-golang/
func Info(a ...any) {
	prefix := time.Now().Format(time.RFC3339)
	message := prefix + addColorString(WHITE, " [INFO] "+fmt.Sprint(a...))
	loggerVar.Println(message)
}

func Warn(a ...any) {
	prefix := time.Now().Format(time.RFC3339)
	message := prefix + addColorString(YELLOW, " [WARN] "+fmt.Sprint(a...))
	loggerVar.Println(message)
}

func Error(a ...any) {
	prefix := time.Now().Format(time.RFC3339)
	message := prefix + addColorString(RED, " [ERROR] "+fmt.Sprint(a...))
	loggerVar.Println(message)
}

func addColorString(color Color, s string) string {
	var stringColor string
	switch color {
	case YELLOW:
		stringColor = "\033[33m"
	case RED:
		stringColor = "\033[31m"
	case WHITE:
		stringColor = ""
	}
	return stringColor + s + "\033[0m"
}
