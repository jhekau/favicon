package checks

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 24 July 2023
 * проверка исходного изображения
 */
import (
	"fmt"

	config_ "github.com/jhekau/favicon/internal/config"
	logger_ "github.com/jhekau/favicon/internal/core/logger"
	storage_ "github.com/jhekau/favicon/pkg/models/storage"
)

const (
	logCP  = `/internal/service/convert/checks/source.go`
	logC01 = `C01: check source`
	logC02 = `C02: check is exists source file`
	logC03 = `C03: file is not exist`
	logC04 = `C04: get resolution source file`
	logC05 = `C05: the resolution of the preview is larger than the original image`
	logC06 = `C06: incorrect resolution source file`
	// logC07 = `C07: `
	// logC08 = `C08: `
	// logC09 = `C09: `
)

//
type Source struct {
	L *logger_.Logger
	Cache interface{
		Status(original storage_.StorageKey, originalSVG bool, thumb_size int) (bool, error)
		SetErr(original storage_.StorageKey, originalSVG bool, thumb_size int, err error) error
	}
	Resolution interface{
		Get(_ storage_.StorageOBJ) (w int, h int, err error)
	}
}

// проверка исходного изображения на корректность
func (c *Source) Check( original storage_.StorageOBJ, originalSVG bool, thumb_size int ) error {

	ok, err := c.Cache.Status(original.Key(), originalSVG, thumb_size)
	if ok {
		if err == nil {
			return nil
		}
		return c.L.Typ.Error(logCP, logC01, err)
	}

	exist, err := original.IsExists()
	if err != nil {
		return c.Cache.SetErr(original.Key(), originalSVG, thumb_size, c.L.Typ.Error(logCP, logC02, err))
	}
	if !exist {
		return c.Cache.SetErr(original.Key(), originalSVG, thumb_size, c.L.Typ.Error(logCP, logC03, err))
	}

	if !originalSVG {

		source_width, source_height, err := c.Resolution.Get(original)
		if err != nil {
			return c.Cache.SetErr(original.Key(), originalSVG, thumb_size, c.L.Typ.Error(logCP, logC04, err))
		}

		if source_height < thumb_size || source_width < thumb_size {
			return c.Cache.SetErr(
				original.Key(), originalSVG, thumb_size, c.L.Typ.Error(logCP, 
					fmt.Sprintf(`Source: %vx%v, Preview: %v; %s`, source_width, source_height, thumb_size, logC05),
					err),
				)
		}

		if 	source_height < config_.ImageSourceResolutionMin || 
			source_height > config_.ImageSourceResolutionMax ||
			source_width < config_.ImageSourceResolutionMin || 
			source_width > config_.ImageSourceResolutionMax {

				return c.Cache.SetErr(
					original.Key(), originalSVG, thumb_size, c.L.Typ.Error(logCP, 
						fmt.Sprintf(`Min Resolution: %v, Max Resolution: %v, Current Resolution: %vx%v`,
						config_.ImageSourceResolutionMin,
						config_.ImageSourceResolutionMax,
						source_width,
						source_height,
					), logC06, err))
		}
	}

	c.Cache.SetErr(original.Key(), originalSVG, thumb_size, nil)
	return nil
}
