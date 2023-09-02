package err

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 16 August 2023
 */
import (
	"fmt"

	"github.com/jhekau/favicon/interfaces/logger"
)

func Err(log logger.Logger, path string, arg ...any) error {
	log.Error(path, arg)
	return fmt.Errorf("[%s] %v ", path, fmt.Sprint(arg...))
}