package convert

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
	logF02 = `F02: check preview image`
	logF03 = `F03: convert thumb`
	logF04 = `F04: check source image`
	logF05 = `F05: the resolution is 0`
	logF06 = `F06: there is no suitable converter for image modification`
)
func errF(i... interface{}) error {
	return err_.Err(err_.TypeError, `/internal/service/convert.go`, i...)
} 

var (
	Convert = convert_file
)



// конвертер, который непосредственно занимается конвертацией
type Converter interface {
	Proc(source, save types_.FilePath, size_px int, typ types_.FileType) error
}

// проверяют тип, в который хотим конвертировать и запускают конвертер
type ConvertType interface {

	// функция, которая запускается для проверки типа и последующего запуска конвертора
	Do(source, save types_.FilePath, size_px int, typ types_.FileType, conv Converter) (complete bool, err error)
}

// связка, функция запускающая конвертер и сам конвертер
type Converters struct {
	Converter Converter
	ConvertType ConvertType
}

// for testing
// var fn_source_check func( fpath types_.FilePath, source_typ types_.FileType, thumb_size int ) error
// func init(){
// 	fn_source_check = source_check
// }

// проверка валидности запрашиваемой превьюхи
type CheckPreview interface {
	Check(typ types_.FileType, size_px int) error
}

// проверка валидности исходника для нарезания превьюхи
type CheckSource interface {
	Check(fpath types_.FilePath, source_typ types_.FileType, thumb_size int) error
}


// конвертирование исходного изображения нужную превьюшку
func convert_file( 
	source, source_svg, save types_.FilePath,
	typ types_.FileType,
	size_px int,
	converters []Converters,
	check_preview CheckPreview,
	check_source CheckSource,
)(
	err error,
){

	// default: условно
	source_type := types_.PNG()

	if size_px == 0 {
		return errF(logF05, err)
	}

	if source == `` {
		if source_svg == `` {
			return errF(logF01, err)
		}
		source = source_svg
		source_type = types_.SVG()
	}

	err = check_preview.Check(typ, size_px)
	if err != nil {
		return errF(logF02, err)
	}

	err = check_source.Check(source, source_type, size_px)
	if err != nil {
		return errF(logF04, err)
	}

	// for _, fn := range []func(s, sv types_.FilePath, sz int, tp types_.FileType, conv Converter) (bool, error) {
	// 	convert_ico,
	// 	convert_png,
	// }
	for _, fn := range converters {
		if ok, err := fn.ConvertType.Do(source, save, size_px, typ, fn.Converter); err != nil {
			return errF(logF03, err)
		} else if ok {
			return nil
		}
	}
	return errF(logF06, err)
}






