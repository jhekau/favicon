package converter_exec_anthonynsimon

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 March 2023
 */

import (
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	ico_ "github.com/Kodeworks/golang-image-ico"
	logger_ "github.com/jhekau/favicon/internal/core/logger"
	types_ "github.com/jhekau/favicon/internal/core/types"
)

const (
	logC   = `/thumb/convert/convert.go`
	logC01 = `C01: open source image`
	logC02 = `C02: unsupported type image`
	logC03 = `C03: save thumb image`
)

type Exec struct{
	l *logger_.Logger
}

func (e *Exec) Proc(source, save types_.FilePath, size_px int, typ types_.FileType) error {

    img, err := imgio.Open(source.String())
    if err != nil {
        return e.l.Typ.Error(logC, logC01, err)
    }
	encoder := imgio.PNGEncoder()
	// if typ == types_.ICO() {
	// 	encoder = imgio.Encoder(
	// 		ico_.Encode,
	// 	)
	// }
	switch typ {
	case types_.ICO():
		encoder = imgio.Encoder(
			ico_.Encode,
		)
	case types_.PNG():
	default:
		return e.l.Typ.Error(logC, logC02, err)
	}
	
	if err = imgio.Save(
		save.String(),
		transform.Resize(img, size_px, size_px, transform.Linear),
		encoder,
	); err != nil {
		return e.l.Typ.Error(logC, logC03, err)
    }
	return nil
}