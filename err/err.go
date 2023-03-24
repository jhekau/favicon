package err

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 24 March 2023
 */
import (
	"fmt"
)

type TypeErr string
const (
	TypeError TypeErr = `error`
	TypeAlert TypeErr = `alert`
	TypeInfo TypeErr = `info`
)

func Err(typ TypeErr, path string, messages ...interface{}) error {
	return fmt.Errorf("[%s] %s %v", typ, path, fmt.Sprint(messages...))
}