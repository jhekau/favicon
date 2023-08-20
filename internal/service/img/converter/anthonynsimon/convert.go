package converter_exec_anthonynsimon

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 March 2023
 */

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	ico_ "github.com/Kodeworks/golang-image-ico"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	err_ "github.com/jhekau/favicon/internal/core/err"
	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"
	types_ "github.com/jhekau/favicon/pkg/core/types"
	"github.com/pressly/goico"
)

const (
	logC   = `internal/service/img/converter/anthonynsimon/convert.go`
	logC01 = `C01: error read source`
	logC02 = `C02: error image decode`
	logC03 = `C03: error open writer storage`
	logC04 = `C04: error transform image`
)


func init() {
	image.RegisterFormat("ico", "\x00\x00\x01\x00", ico.Decode, ico.DecodeConfig)
}

type Exec struct{
	L logger_.Logger
}

func (e *Exec) Proc(original, save storage_.StorageOBJ , size_px int, thumbTyp types_.FileType) error {

	r, err := original.Reader()
	if err != nil {
		return err_.Err(e.L, logC, logC01, err)
	}
	
	img, _, err := image.Decode(r)
	if err != nil {
		return err_.Err(e.L, logC, logC02, err)
	}

	encoder := imgio.PNGEncoder()
	if thumbTyp == types_.ICO() {
		encoder = imgio.Encoder(ico_.Encode)
	}

	w, err := save.Writer()
	if err != nil {
		return err_.Err(e.L, logC, logC03, err)
	}
	defer w.Close()

	// s := stor{ obj: &bytes.Buffer{} }
	// w := s.Writer()

	err = encoder(w, transform.Resize(img, size_px, size_px, transform.Linear))
	if err != nil {
		return err_.Err(e.L, logC, logC04, err)
	}

// fmt.Println(` >>>> DEBUG >>>> `, logC, fmt.Sprint( io.ReadAll(s.obj) ))

	return nil
}


// type stor struct {
// 	obj *bytes.Buffer
// }

// func (s *stor) Writer() io.Writer {
// 	return s.obj
// }
