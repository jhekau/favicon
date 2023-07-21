package create

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 21 July 2023
 * backdoor for test
 */
import (
	"os"

	types_ "github.com/jhekau/favicon/internal/core/types"
)

func ConvertFileTest( 
	source, source_svg, save types_.FilePath,
	typ types_.FileType,
	size_px int,
	conv Converter,
)(
	complete bool,
	err error,
){
	fn_source_check = func( fpath types_.FilePath, source_typ types_.FileType, thumb_size int ) error {
		return nil
	}

	return convert_file(source, source_svg, save, typ, size_px, conv)
}

var ConvertICOTest = convert_ico
var ConvertPNGTest = convert_png

func SourceResolutionTest( fpath types_.FilePath, FnOsOpen func(name string) (*os.File, error) ) ( w,h int, err error ) {
	osOpen = FnOsOpen
	return source_resolution( fpath )
}

func SourceCheckTest( fpath types_.FilePath, source_typ types_.FileType, thumb_size int, FnOsStat func(name string) (os.FileInfo, error) ) error {
	osStat = FnOsStat
	return source_check( fpath, source_typ, thumb_size )
}


