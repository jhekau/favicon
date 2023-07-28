package files

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	// "image"
	"io"
	"os"

	err_ "github.com/jhekau/favicon/internal/core/err"
	types_ "github.com/jhekau/favicon/internal/core/types"
)

const (
	// logS01 = `S01: open source image`
	// logS02 = `S02: decode image config`
	logS03 = `S03: os stat suorce image`
	logS04 = `S04: save thumb image is a folder`
)
func errS(i... interface{}) error {
	return err_.Err(err_.TypeError, `/internal/data/files/stat.go`, i...)
} 

var (
	// Resolution = resolution
	Read = read
	IsExists = exists
)

// for test 
var (
	osOpen = os.Open
	osStat = os.Stat
)
/*
// получние разрешения исходного изображения
func resolution( fpath types_.FilePath ) ( w,h int, err error ) {

	file, err := osOpen(fpath.String())
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
*/

func read(fpath types_.FilePath) (io.ReadCloser, error) {
	return osOpen(fpath.String())
}

func exists( fpath types_.FilePath ) ( bool, error ) {

	if f, err := osStat(fpath.String()); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errS(logS03, err)

	} else if f.IsDir() {
		return false, errS(logS04, err)
	}

	return true, nil
}


