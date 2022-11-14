package mylogger

import (
	"github.com/fatih/color"
	"io"
	"log"
)

type MyLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

var logger *MyLogger

func Init(infoWriter, errorWriter io.Writer) {
	logger = &MyLogger{
		infoLogger:  log.New(infoWriter, color.CyanString("[Info] "), log.Ltime|log.Ldate),
		errorLogger: log.New(errorWriter, color.RedString("[Error] "), log.Ltime|log.Ldate),
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
