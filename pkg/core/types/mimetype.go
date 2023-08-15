package types

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 11 August 2023
 */
import (
	"mime"
)

//
type FileType string 
func (f FileType) String() string {
	return string(f)
}
var (
	ico = FileType(mime.TypeByExtension(".ico"))
	svg = FileType(mime.TypeByExtension(".svg"))
	png = FileType(mime.TypeByExtension(".png"))
	jpg = FileType(mime.TypeByExtension(".png"))

	ICO = func() FileType { return ico }
	SVG = func() FileType { return svg }
	PNG = func() FileType { return png }
	JPG = func() FileType { return jpg }
) 