package converter_exec_anthonynsimon

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 March 2023
 */

import (
	"io"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	ico_ "github.com/Kodeworks/golang-image-ico"
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
	l *logger_.Logger
}

func (e *Exec) Proc(source, save StorageOBJ, size_px int, thumbTyp types_.FileType) error {

	r, err := source.Reader()
	if err != nil {
		return e.l.Typ.Error(logC, logC01, err)
	}

	img, _, err := image.Decode(r)
	if err != nil {
		return e.l.Typ.Error(logC, logC02, err)
	}

	encoder := imgio.PNGEncoder()
	if thumbTyp == types_.ICO() {
		encoder = imgio.Encoder(ico_.Encode)
	}

	w, err := save.Writer()
	if err != nil {
		return e.l.Typ.Error(logC, logC03, err)
	}

	err = encoder(w, transform.Resize(img, size_px, size_px, transform.Linear))
	if err != nil {
		return e.l.Typ.Error(logC, logC04, err)
	}

	return nil
}