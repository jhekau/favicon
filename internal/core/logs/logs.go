package logs

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 1 August 2023
 */
import (
	"fmt"
	"log"
)

type Logger struct{}

func (l *Logger) Info(path string, arg ...any) {
	log.Println(fmt.Sprintf("[%s] %s %v: ", `info`, path, fmt.Sprint(arg...)))
}
func (l *Logger) Warn(path string, arg ...any) {
	log.Println(fmt.Sprintf("[%s] %s %v: ", `alert`, path, fmt.Sprint(arg...)))
}
func (l *Logger) Error(path string, arg ...any) {
	log.Println(fmt.Sprintf("[%s] %s %v: ", `error`, path, fmt.Sprint(arg...)))
}

