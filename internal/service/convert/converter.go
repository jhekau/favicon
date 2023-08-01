package convert

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 26 August 2023
 */
import(
	err_ "github.com/jhekau/favicon/internal/core/err"
	types_ "github.com/jhekau/favicon/internal/core/types"
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

// конвертер, который непосредственно занимается конвертацией
/*
type ConverterExec interface{
	Proc(source, save types_.FilePath, size_px int, typ types_.FileType) error
}

// проверяют тип, в который хотим конвертировать и запускают конвертер
// метод, которая запускается для проверки типа и последующего запуска конвертора
type ConverterType interface{
	Do(source, save types_.FilePath, size_px int, typ types_.FileType, conv ConverterExec) (complete bool, err error)
}

// связка, функция запускающая конвертер и сам конвертер
type Converters struct {
	Converter ConverterExec
	ConvertType ConverterType
}
*/
// конвертер для конкретного типа
type ConverterT interface{
	Do(source, save types_.FilePath, size_px int, typ types_.FileType) (complete bool, err error)
}

// проверка валидности запрашиваемой превьюхи
type CheckPreview interface {
	Check(typ types_.FileType, size_px int) error
}

// проверка валидности исходника для нарезания превьюхи
type CheckSource interface {
	Check(fpath types_.FilePath, source_typ types_.FileType, thumb_size int) error
}




// конвертирование исходного изображения нужную превьюшку
type Converter struct{
	Converters []ConverterT
	CheckPreview CheckPreview
	CheckSource CheckSource
}

func (c *Converter) Do( 
	source, source_svg, save types_.FilePath,
	typ types_.FileType,
	size_px int,
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

	err = c.CheckPreview.Check(typ, size_px)
	if err != nil {
		return errF(logF02, err)
	}

	err = c.CheckSource.Check(source, source_type, size_px)
	if err != nil {
		return errF(logF04, err)
	}

	for _, fn := range c.Converters {
		// if ok, err := fn.ConvertType.Do(source, save, size_px, typ, fn.Converter); err != nil {
		if ok, err := fn.Do(source, save, size_px, typ); err != nil {
			return errF(logF03, err)
		} else if ok {
			return nil
		}
	}
	return errF(logF06, err)
}


