package convert

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 March 2023
 */

import (
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	ico_ "github.com/Kodeworks/golang-image-ico"

	types_ "github.com/jhekau/favicon/types"
)

var (
	// ~~ interface ~~
	Do = do
)

func do(source, save types_.FilePath, size_px int, typ types_.FileType) error {

    img, err := imgio.Open(source.String())
    if err != nil {
        // return error
    }
	encoder := imgio.PNGEncoder()
	if typ == types_.ICO() {
		encoder = imgio.Encoder(
			ico_.Encode,
		)
	}
	
	if err = imgio.Save(
		save.String(),
		transform.Resize(img, size_px, size_px, transform.Linear),
		encoder,
	); err != nil {
        // return error
    }
	return nil
}