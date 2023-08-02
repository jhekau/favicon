package convert

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 1 August 2023
 */
import(
	logger_ "github.com/jhekau/favicon/internal/core/logger"
	types_ "github.com/jhekau/favicon/internal/core/types"
)

const (
	logFP  = `/internal/service/convert.go`
	logF01 = `F01: empty original image`
	logF02 = `F02: check preview image`
	logF03 = `F03: convert thumb`
	logF04 = `F04: check source image`
	logF05 = `F05: the resolution is 0`
	logF06 = `F06: there is no suitable converter for image modification`
)

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
	Check(fpath types_.FilePath, originalSVG bool, thumb_size int) error
}




// конвертирование исходного изображения нужную превьюшку
type Converter struct{
	L *logger_.Logger
	Converters []ConverterT
	CheckPreview CheckPreview
	CheckSource CheckSource
}

func (c *Converter) Do( 
	source, /*source_svg,*/ save types_.FilePath,
	originalSVG bool,
	typThumb types_.FileType,
	size_px int,
)(
	err error,
){

	if size_px == 0 {
		return c.L.Typ.Error(logFP, logF05, err)
	}
	if source == `` {
		return c.L.Typ.Error(logFP, logF01, err)
	}

	err = c.CheckPreview.Check(typThumb, size_px)
	if err != nil {
		return c.L.Typ.Error(logFP, logF02, err)
	}

	err = c.CheckSource.Check(source, originalSVG, size_px)
	if err != nil {
		return c.L.Typ.Error(logFP, logF04, err)
	}

	for _, fn := range c.Converters {
		// if ok, err := fn.ConvertType.Do(source, save, size_px, typ, fn.Converter); err != nil {
		if ok, err := fn.Do(source, save, size_px, typThumb); err != nil {
			return c.L.Typ.Error(logFP, logF03, err)
		} else if ok {
			return nil
		}
	}
	return c.L.Typ.Error(logFP, logF06, err)
}


