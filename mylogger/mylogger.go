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

func New() *MyLogger {
	logger = &MyLogger{
		infoLogger:  log.New(os.Stdout, color.CyanString("[Info] "), log.Ltime|log.Ldate),
		errorLogger: log.New(os.Stderr, color.RedString("[Error] "), log.Ltime|log.Ldate),
	}
	return logger
}

func Get() *MyLogger {
	return logger
}

func (l *MyLogger) Println(v ...any) {
	l.infoLogger.Println(v...)
}

func (l *MyLogger) Errorln(v ...any) {
	l.errorLogger.Println(v...)
}

func (l *MyLogger) Fatalln(v ...any) {
	l.errorLogger.Fatalln(v...)
}
