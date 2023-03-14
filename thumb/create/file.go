package create

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	types_ "github.com/jhekau/favicon/types"
	converter_ "github.com/jhekau/favicon/thumb/convert"
)

var (
	// ~~ interface ~~

	File = convert_file
)


//
func convert_file( source, source_svg, save types_.FilePath, typ types_.FileType, size_px int ) (complete bool, err error) {

	// default: условно
	source_type := types_.PNG()

	if source == `` {
		if source_svg == `` {
			// return error
		}
		source = source_svg
		source_type = types_.SVG()
	}

	err = source_check(source, source_type, size_px)
	if err != nil {
		// return error
	}

	for _, fn := range []func(s,sv types_.FilePath, sz int, tp types_.FileType) (bool, error) {
		convert_ico,
		convert_png,
	}{
		if ok, err := fn(source, save, size_px, typ); err != nil {
			// return error
		} else if ok {
			return true, nil
		}
	}
	return false, nil
}

//
func convert_ico(source, save types_.FilePath, size_px int, typ types_.FileType) (complete bool, err error) {
	if typ != types_.ICO() {
		return false, nil
	}
    if err := converter_.Do(source, save, size_px, typ); err != nil {
		// return error
	}
	return true, nil
}

//
func convert_png(source, save types_.FilePath, size_px int, typ types_.FileType) (complete bool, err error) {
	if typ != types_.PNG() {
		return false, nil
	}
    if err := converter_.Do(source, save, size_px, typ); err != nil {
		// return error
	}
	return true, nil
}

