package resolution

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 28 July 2023
 * получение разрешение изображения
 */
import (
	"image"
	"io"

	err_ "github.com/jhekau/favicon/internal/core/err"
)

const (
	logR01 = `R01: image decode config error`
	// logR02 = `R02: `
	// logR03 = `R03: `
	// logR04 = `R04: `
)
func errR(i... interface{}) error {
	return err_.Err(err_.TypeError, `/favicon/internal/service/img/resolution/resolution.go`, i...)
}

// Get : получние разрешения изображения 
func Get(r io.Reader) ( w,h int, err error ) {
	image, _, err := image.DecodeConfig(r)
    if err != nil {
		return 0,0, errR(logR01, err)
    }
    return image.Width, image.Height, nil
}