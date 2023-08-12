package logger_mock

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 1 August 2023
 * from test
 * Предоставляет возможноть протестировать возвращаемые ошибки юнитами
 */
import (
	"fmt"
)

type LoggerErrorf struct{}

func (l *LoggerErrorf) Info(path string, messages ...interface{}) error {
	return fmt.Errorf("[%s] %s %v", `info`, path, fmt.Sprint(messages...))
}
func (l *LoggerErrorf) Alert(path string, messages ...interface{}) error {
	return fmt.Errorf("[%s] %s %v", `info`, path, fmt.Sprint(messages...))
}
func (l *LoggerErrorf) Error(path string, messages ...interface{}) error {
	return fmt.Errorf("[%s] %s %v", `info`, path, fmt.Sprint(messages...))
}
