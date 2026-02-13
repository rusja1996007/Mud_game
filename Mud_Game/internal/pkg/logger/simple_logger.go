package logger

import (
	"fmt"
	"time"
)

type SimpleLogger struct {
	name string
}

// log - внутренний метод для форматирования и вывода
func (l *SimpleLogger) log(level, msg string) {
	timeStamp := time.Now().Format("2006-01-02 15:15:15")
	fmt.Println(timeStamp, level, l.name, msg)
}

func NewSimpleLogger(name string) *SimpleLogger {
	return &SimpleLogger{name: name}
}
func (l *SimpleLogger) Info(msg string) {
	l.log("INFO", msg)
}
func (l *SimpleLogger) Error(msg string) {
	l.log("ERROR", msg)
}

func (l *SimpleLogger) Debug(msg string) {
	l.log("DEBUG", msg)
}

func (l *SimpleLogger) Warn(msg string) {
	l.log("WARN", msg)
}
