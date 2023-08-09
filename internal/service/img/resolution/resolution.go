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
	logR02 = `R02: read object`
	// logR03 = `R03: `
	// logR04 = `R04: `
)

type StorageOBJ interface{
	Reader() (io.ReadCloser , error)
}

type Resolution struct {
	L *logger_.Logger
}

// Get : получние разрешения изображения 
func (r *Resolution) Get(obj StorageOBJ) ( w,h int, err error ) {

	read, err := obj.Reader()
	if err != nil {
		return 0,0, r.L.Typ.Error(logRP, logR02, err)
	}
	defer read.Close()

	image, _, err := image.DecodeConfig(read)
    if err != nil {
		return 0,0, r.L.Typ.Error(logRP, logR01, err)
    }
    return image.Width, image.Height, nil
}
