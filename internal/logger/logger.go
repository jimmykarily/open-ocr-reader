// Package logger implement a Logger class for this project
package logger

import (
	"log"
	"os"
)

type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

func New() Logger {
	return Logger{
		InfoLogger:  log.New(os.Stdout, "", log.LstdFlags),
		ErrorLogger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (l Logger) Log(msg string) {
	l.InfoLogger.Println(msg)
}

func (l Logger) Logf(msg string, v ...any) {
	l.InfoLogger.Printf(msg, v)
}

func (l Logger) Error(msg string) {
	l.ErrorLogger.Println(msg)
}

func (l Logger) Errorf(msg string, v ...any) {
	l.ErrorLogger.Printf(msg, v)
}
