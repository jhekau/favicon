package logger

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 August 2023
 */
import()

type Logger interface {
	Info(path string, messages ...interface{}) error
	Alert(path string, messages ...interface{}) error
	Error(path string, messages ...interface{}) error
}