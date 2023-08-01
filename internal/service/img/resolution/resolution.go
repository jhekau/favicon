package resolution

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 28 July 2023
 * получение разрешение изображения
 */
import (
	"image"
	"io"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
)

const (
	logRP = `/favicon/internal/service/img/resolution/resolution.go`
	logR01 = `R01: image decode config error`
	// logR02 = `R02: `
	// logR03 = `R03: `
	// logR04 = `R04: `
)

type Resolution struct {
	L *logger_.Logger
	Reader interface{
		Read() io.Reader
	}
}

// Get : получние разрешения изображения 
func (r *Resolution) Get() ( w,h int, err error ) {
	image, _, err := image.DecodeConfig(r.Reader.Read())
    if err != nil {
		return 0,0, r.L.Typ.Error(logRP, logR01, err)
    }
    return image.Width, image.Height, nil
}
