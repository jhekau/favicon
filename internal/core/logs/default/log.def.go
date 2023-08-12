package logger_default

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 1 August 2023
 */
import (
	"fmt"
)

type Logger struct{}

func (l *Logger) Info(path string, messages ...interface{}) error {
	return fmt.Errorf("[%s] %s %v", `info`, path, fmt.Sprint(messages...))
}
func (l *Logger) Alert(path string, messages ...interface{}) error {
	return fmt.Errorf("[%s] %s %v", `alert`, path, fmt.Sprint(messages...))
}
func (l *Logger) Error(path string, messages ...interface{}) error {
	return fmt.Errorf("[%s] %s %v", `error`, path, fmt.Sprint(messages...))
}

