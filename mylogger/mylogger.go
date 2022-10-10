package mylogger

import (
	"github.com/fatih/color"
	"log"
	"os"
)

type MyLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

var logger *MyLogger

func init() {
	logger = &MyLogger{
		infoLogger:  log.New(os.Stdout, color.CyanString("[Info] "), log.Ltime|log.Ldate),
		errorLogger: log.New(os.Stderr, color.RedString("[Error] "), log.Ltime|log.Ldate),
	}
}

func Println(v ...any) {
	logger.infoLogger.Println(v...)
}

func Errorln(v ...any) {
	logger.errorLogger.Println(v...)
}

func Fatalln(v ...any) {
	logger.errorLogger.Fatalln(v...)
}
