package logger_mock

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 1 August 2023
 * from test
 */
import ()

type LoggerErrorSet struct{
	Err error
}

func (l *LoggerErrorSet) Info(path string, messages ...interface{}) error {
	return l.Err
}
func (l *LoggerErrorSet) Alert(path string, messages ...interface{}) error {
	return l.Err
}
func (l *LoggerErrorSet) Error(path string, messages ...interface{}) error {
	return l.Err
}


