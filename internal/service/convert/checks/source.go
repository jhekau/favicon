package checks

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 24 July 2023
 * проверка исходного изображения
 */
import (
	"fmt"
	"sync"

	config_ "github.com/jhekau/favicon/internal/config"
	err_ "github.com/jhekau/favicon/internal/core/err"
	types_ "github.com/jhekau/favicon/internal/core/types"
	files_ "github.com/jhekau/favicon/internal/data/files"
)

const (
	logC01 = `C01: check source`
	logC02 = `C02: check is exists source file`
	logC03 = `C03: file is not exist`
	logC04 = `C04: get resolution source file`
	logC05 = `C05: the resolution of the preview is larger than the original image`
	logC06 = `C06: incorrect resolution source file`
	logC07 = `C07: `
	logC08 = `C08: `
	logC09 = `C09: `
)
func errC(i... interface{}) error {
	return err_.Err(err_.TypeError, `/internal/service/convert/checks/source.go`, i...)
} 

// список оригинальных файлов, ранее проверенных на корректность
var check = struct{
	Status func(fpath types_.FilePath) (bool, error)
	SetErr func(fpath types_.FilePath, err error) error
}{}

func init() {

	var c sync.Map

	check.Status = func(p types_.FilePath) (bool, error) {
		e, ok := c.Load(p)
		return ok, e.(error)
	}

	check.SetErr = func(fpath types_.FilePath, err error) error {
		c.Store(fpath, err)
		return err
	}
}

//
type Source struct {}

// проверка исходного изображения на корректность
func (c Source) Check( fpath types_.FilePath, source_typ types_.FileType, thumb_size int ) error {

	ok, err := check.Status(fpath)
	if ok {
		if err == nil {
			return nil
		}
		return errC(logC01, err)
	}

	exist, err := files_.IsExists(fpath)
	if err != nil {
		return check.SetErr(fpath, errC(logC02, err))
	}
	if !exist {
		return check.SetErr(fpath, errC(logC03, err))
	}

	if source_typ != types_.SVG() {

		source_width, source_height, err := files_.Resolution(fpath)
		if err != nil {
			return check.SetErr(fpath, errC(logC04, err))
		}

		if source_height < thumb_size || source_width < thumb_size {
			return check.SetErr(
				fpath, errC(
					fmt.Sprintf(`Source: %vx%v, Preview: %v; %s`, source_width, source_height, thumb_size, logC05),
					err),
				)
		}

		if 	source_height < config_.ImageSourceResolutionMin || 
			source_height > config_.ImageSourceResolutionMax ||
			source_width < config_.ImageSourceResolutionMin || 
			source_width > config_.ImageSourceResolutionMax {

				return check.SetErr(
					fpath, errC(
						fmt.Sprintf(`Min Resolution: %v, Max Resolution: %v, Current Resolution: %vx%v`,
						config_.ImageSourceResolutionMin,
						config_.ImageSourceResolutionMax,
						source_width,
						source_height,
					), logC06, err))
		}
	}

	check.SetErr(fpath, nil)
	return nil
}
