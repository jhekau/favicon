package create

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	err_ "github.com/jhekau/favicon/internal/core/err"
	types_ "github.com/jhekau/favicon/internal/core/types"
	// converter_ "github.com/jhekau/favicon/thumb/convert"
)


const (
	logF01 = `F01: thumb is SVG, false convert`
	logF02 = `F02: check source image`
	logF03 = `F03: convert thumb`
	logF04 = `F04: convert ico`
	logF05 = `F05: convert png`
)
func errF(i... interface{}) error {
	return err_.Err(err_.TypeError, `/thumb/create/file.go`, i...)
} 

var (
	File = convert_file
)

type Converter interface {
	Do(source, save types_.FilePath, size_px int, typ types_.FileType) error
}

//
func convert_file( 
	source, source_svg, save types_.FilePath,
	typ types_.FileType,
	size_px int,
	conv Converter,
)(
	complete bool,
	err error,
){

	// default: условно
	source_type := types_.PNG()

	if size_px == 0 {
		return false, nil
	}

	if source == `` {
		if source_svg == `` {
			return false, errF(logF01, err)
		}
		source = source_svg
		source_type = types_.SVG()
	}

	err = source_check(source, source_type, size_px)
	if err != nil {
		return false, errF(logF02, err)
	}

	for _, fn := range []func(s, sv types_.FilePath, sz int, tp types_.FileType, conv Converter) (bool, error) {
		convert_ico,
		convert_png,
	}{
		if ok, err := fn(source, save, size_px, typ, conv); err != nil {
			return false, errF(logF03, err)
		} else if ok {
			return true, nil
		}
	}
	return false, nil
}

//
func convert_ico(source, save types_.FilePath, size_px int, typ types_.FileType, conv Converter) (complete bool, err error) {
	if typ != types_.ICO() {
		return false, nil
	}
    if err := conv.Do(source, save, size_px, typ); err != nil {
		return false, errF(logF04, err)
	}
	return true, nil
}

//
func convert_png(source, save types_.FilePath, size_px int, typ types_.FileType, conv Converter) (complete bool, err error) {
	if typ != types_.PNG() {
		return false, nil
	}
    if err := conv.Do(source, save, size_px, typ); err != nil {
		return false, errF(logF05, err)
	}
	return true, nil
}

