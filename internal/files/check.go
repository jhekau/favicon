package create

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	"errors"
	"image"
	"os"
	"sync"

	err_ "github.com/jhekau/favicon/core/err"
	types_ "github.com/jhekau/favicon/core/types"
)

const (
	logS01 = `S01: open source image`
	logS02 = `S02: decode image config`
	logS03 = `S03: os stat suorce image`
	logS04 = `S04: save thumb image is a folder`
	logS05 = `S05: get source resolution`
	logS06 = `S06: source size < thumb size`
)
func errS(i... interface{}) error {
	return err_.Err(err_.TypeError, `/thumb/create/source.go`, i)
} 

// список оригинальных файлов, ранее проверенных на корректность
var source_check_list sync.Map // [filepath]error

// получние разрешения исходного изображения
func source_resolution( fpath types_.FilePath ) ( w,h int, err error ) {

	file, err := os.Open(fpath.String())
    if err != nil {
		return 0,0, errS(logS01, err)
    }
	defer file.Close()

    image, _, err := image.DecodeConfig(file)
    if err != nil {
		return 0,0, errS(logS02, err)
    }
    return image.Width, image.Height, nil
}

var errOK = errors.New(`OK`)

// проверка исходного изображения на корректность
func source_check( fpath types_.FilePath, source_typ types_.FileType, thumb_size int ) error {

	if e, ok := source_check_list.Load(fpath); ok {
		err := e.(error)
		if err == errOK {
			return nil
		}
		return e.(error)
	}

	if f, err := os.Stat(fpath.String()); err != nil {

		source_check_list.Store(fpath, err)
		return errS(logS03, err)

	} else if f.IsDir() {
		err := errS(logS04, err)
		source_check_list.Store(fpath, err)
		return err
	}

	if source_typ != types_.SVG() {

		source_width, source_height, err := source_resolution(fpath)
		if err != nil {
			err := errS(logS05, err)
			source_check_list.Store(fpath, err)
			return err
		}

		if  source_height < thumb_size || source_width < thumb_size {
			err := errS(logS06, err)
			source_check_list.Store(fpath, err)
			return err
		}
	}

	source_check_list.Store(fpath, errOK)
	return nil
}
