package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)}
}

func (l *Logger) Error(err error) {
	l.Printf("ERROR: %s", err)
}
