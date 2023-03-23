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

	types_ "github.com/jhekau/favicon/types"
)

// список оригинальных файлов, ранее проверенных на корректность
var source_check_list sync.Map // [filepath]error

// ...
func source_resolution( fpath types_.FilePath ) ( w,h int, err error ) {

	file, err := os.Open(fpath.String())
    if err != nil {
		// return error
    }
	defer file.Close()

    image, _, err := image.DecodeConfig(file)
    if err != nil {
		// return error
    }
    return image.Width, image.Height, nil
}

var errOK = errors.New(`OK`)

// ...
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
		// return error
	} else if f.IsDir() {
		// new error - конечное превью является каталогом
		// source_check_list.Store(fpath, err)
		// return error
	}

	source_width, source_height, err := source_resolution(fpath)
	if err != nil {
		// source_check_list.Store(fpath, err)
		// return error
	}

	if source_typ != types_.SVG() && ( source_height < thumb_size || source_width < thumb_size ) {
		// source_check_list.Store(fpath, err)
		// return error
	}

	source_check_list.Store(fpath, errOK)
	return nil
}
