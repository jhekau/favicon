package convert

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 1 August 2023
 */
import (
	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	types_ "github.com/jhekau/favicon/pkg/core/types"
	converter_ "github.com/jhekau/favicon/pkg/core/models/converter"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"
	err_ "github.com/jhekau/favicon/internal/core/err"
)

const (
	logFP  = `/internal/service/convert.go`
	// logF01 = `F01: `
	logF02 = `F02: check preview image`
	logF03 = `F03: convert thumb`
	logF04 = `F04: check source image`
	logF05 = `F05: the resolution is 0`
	logF06 = `F06: there is no suitable converter for image modification`
)

// проверка валидности запрашиваемой превьюхи
type CheckPreview interface {
	Check(typ types_.FileType, size_px int) error
}

// проверка валидности исходника для нарезания превьюхи
type CheckSource interface {
	Check(original storage_.StorageOBJ, originalSVG bool, thumb_size int) error
}





// конвертирование исходного изображения нужную превьюшку
type Converter struct{
	L logger_.Logger
	Converters []converter_.ConverterTyp
	CheckPreview CheckPreview
	CheckSource CheckSource
}

func (c *Converter) Do( 
	source, /*source_svg,*/ save storage_.StorageOBJ, // types_.FilePath,
	originalSVG bool,
	typThumb types_.FileType,
	size_px int,
)(
	err error,
){

	if size_px <= 0 {
		return err_.Err(c.L, logFP, logF05)
	}

	err = c.CheckPreview.Check(typThumb, size_px)
	if err != nil {
		return err_.Err(c.L, logFP, logF02, err)
	}

	err = c.CheckSource.Check(source, originalSVG, size_px)
	if err != nil {
		return err_.Err(c.L, logFP, logF04, err)
	}

	for _, fn := range c.Converters {
		if ok, err := fn.Do(source, save, size_px, typThumb); err != nil {
			return err_.Err(c.L, logFP, logF03, err)
		} else if ok {
			return nil
		}
	}
	return err_.Err(c.L, logFP, logF06)
}


