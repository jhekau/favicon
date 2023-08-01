package logger

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 1 August 2023
 */
import ()


type Logger struct{
	Typ interface {
		Info(path string, messages ...interface{}) error
		Alert(path string, messages ...interface{}) error
		Error(path string, messages ...interface{}) error
	}
}

