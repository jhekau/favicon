package types

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 09 March 2023
 * TDD
 */

import (
	"mime"
)

//
type FileName string
func (f FileName) String() string {
	return string(f)
}

type FilePath string
func (f FilePath) String() string {
	return string(f)
}

type Folder string
func (f Folder) String() string {
	return string(f)
}

//
type FileType string 
func (f FileType) String() string {
	return string(f)
}
var (
	ico = FileType(mime.TypeByExtension(".ico"))
	svg = FileType(mime.TypeByExtension(".svg"))
	png = FileType(mime.TypeByExtension(".png"))

	ICO = func() FileType { return ico }
	SVG = func() FileType { return svg }
	PNG = func() FileType { return png }
) 

//
type FileExists int
const (
	FileExistsNotCheck FileExists = 0
	FileExistsOK FileExists = 1
	FileExistsNOT FileExists = -1
)



//
type URLName string
func (n URLName) String() string {
	return string(n)
}