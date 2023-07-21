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
	Convert = convert_file
	ConverterICO = convert_ico
	ConverterPNG = convert_png
)



// функция, которая непосредственно конвертирует изображение.
// Можно использовать внешнюю библиотеку, или внешний бинарник
type Converter interface {
	Proc(source, save types_.FilePath, size_px int, typ types_.FileType) error
}

// функция, которая проверяет, подходит ли она для конвертации запрашиваемого типа
// и запускает конвертер для конвертации непосредственно файла
type ConvertType interface {

	// функция, которая запускается для проверки типа и последующего запуска конвертора
	Do(source, save types_.FilePath, size_px int, typ types_.FileType, conv Converter) (complete bool, err error)
}

// соб-но пара, интерфейс для конвертера и сам конвертер
type Converters struct {
	Converter Converter
	ConvertType ConvertType
}

// for testing
var fn_source_check func( fpath types_.FilePath, source_typ types_.FileType, thumb_size int ) error
func init(){
	fn_source_check = source_check
}



// конвертирование исходного изображения нужную превьюшку
func convert_file( 
	source, source_svg, save types_.FilePath,
	typ types_.FileType,
	size_px int,
	converters []Converters,
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

	err = fn_source_check(source, source_type, size_px)
	if err != nil {
		return false, errF(logF02, err)
	}

	// for _, fn := range []func(s, sv types_.FilePath, sz int, tp types_.FileType, conv Converter) (bool, error) {
	// 	convert_ico,
	// 	convert_png,
	// }
	for _, fn := range converters {
		if ok, err := fn.ConvertType.Do(source, save, size_px, typ, fn.Converter); err != nil {
			return false, errF(logF03, err)
		} else if ok {
			return true, nil
		}
	}
	return false, nil
}




// интерфейсы для конвертеров

// интерфейс для конвертора ICO
func convert_ico(source, save types_.FilePath, size_px int, typ types_.FileType, conv Converter) (complete bool, err error) {
	if typ != types_.ICO() {
		return false, nil
	}
    if err := conv.Proc(source, save, size_px, typ); err != nil {
		return false, errF(logF04, err)
	}
	return true, nil
}

// интерфейс для конвертора PNG
func convert_png(source, save types_.FilePath, size_px int, typ types_.FileType, conv Converter) (complete bool, err error) {
	if typ != types_.PNG() {
		return false, nil
	}
    if err := conv.Proc(source, save, size_px, typ); err != nil {
		return false, errF(logF05, err)
	}
	return true, nil
}

