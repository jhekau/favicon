package converter_exec_anthonynsimon

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 March 2023
 */

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	ico_ "github.com/Kodeworks/golang-image-ico"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	logger_ "github.com/jhekau/favicon/internal/core/logger"
	types_ "github.com/jhekau/favicon/internal/core/types"
)

const (
	logC   = `internal/service/img/converter/anthonynsimon/convert.go`
	logC01 = `C01: error read source`
	logC02 = `C02: error image decode`
	logC03 = `C03: error open writer storage`
	logC04 = `C04: error transform image`
)


type StorageOBJ interface{
	Reader() (io.ReadCloser, error)
	Writer() (io.WriteCloser, error)
}

type Exec struct{
	L *logger_.Logger
}

func (e *Exec) Proc(original, save StorageOBJ, size_px int, thumbTyp types_.FileType) error {

	r, err := original.Reader()
	if err != nil {
		return e.L.Typ.Error(logC, logC01, err)
	}
	
	img, _, err := image.Decode(r)
	if err != nil {
		return e.L.Typ.Error(logC, logC02, err)
	}

	encoder := imgio.PNGEncoder()
	if thumbTyp == types_.ICO() {
		encoder = imgio.Encoder(ico_.Encode)
	}

	w, err := save.Writer()
	if err != nil {
		return e.L.Typ.Error(logC, logC03, err)
	}
	defer w.Close()

	err = encoder(w, transform.Resize(img, size_px, size_px, transform.Linear))
	if err != nil {
		return e.L.Typ.Error(logC, logC04, err)
	}

	return nil
}