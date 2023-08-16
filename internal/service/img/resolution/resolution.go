package resolution

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 28 July 2023
 * получение разрешение изображения
 */
import (
	"image"

	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"
	err_ "github.com/jhekau/favicon/internal/core/err"
)

const (
	logRP = `/favicon/internal/service/img/resolution/resolution.go`
	logR01 = `R01: image decode config error`
	logR02 = `R02: read object`
	// logR03 = `R03: `
	// logR04 = `R04: `
)

type Resolution struct {
	L logger_.Logger
}

// Get : получние разрешения изображения 
func (r *Resolution) Get(obj storage_.StorageOBJ) ( w,h int, err error ) {

	read, err := obj.Reader()
	if err != nil {
		return 0,0, err_.Err(r.L, logRP, logR02, err)
	}
	defer read.Close()

	image, _, err := image.DecodeConfig(read)
    if err != nil {
		return 0,0, err_.Err(r.L, logRP, logR01, err)
    }
    return image.Width, image.Height, nil
}
