package checks

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 24 July 2023
 * проверка исходного изображения
 */
import (
	"fmt"
	"sync"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	config_ "github.com/jhekau/favicon/internal/config"
	types_ "github.com/jhekau/favicon/internal/core/types"
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

type CacheStatus struct {
	m sync.Map
}
func (c *CacheStatus) cache_key(fpath types_.FilePath, originalSVG bool, thumb_size int) string {
	return fmt.Sprintf(`%s, %v, %d`, fpath, originalSVG, thumb_size)
}
func (c *CacheStatus) Status(fpath types_.FilePath, originalSVG bool, thumb_size int) (bool, error) {

	e, ok := c.m.Load(c.cache_key(fpath, originalSVG, thumb_size))

	var err error
	if e != nil {
		err = e.(error)
	}

	return ok, err
}
func (c *CacheStatus) SetErr(fpath types_.FilePath, originalSVG bool, thumb_size int, err error) error {
	c.m.Store(c.cache_key(fpath, originalSVG, thumb_size), err)
	return err
}

//
type Source struct {
	L *logger_.Logger
	Cache interface{
		Status(fpath types_.FilePath, originalSVG bool, thumb_size int) (bool, error)
		SetErr(fpath types_.FilePath, originalSVG bool, thumb_size int, err error) error
	}
	FileIsExist func(fpath types_.FilePath, l *logger_.Logger ) (bool, error) 
	Resolution interface{
		Get() (w int, h int, err error)
	}
}

// проверка исходного изображения на корректность
func (c *Source) Check( fpath types_.FilePath, originalSVG bool, thumb_size int ) error {

	ok, err := c.Cache.Status(fpath, originalSVG, thumb_size)
	if ok {
		if err == nil {
			return nil
		}
		return c.L.Typ.Error(logCP, logC01, err)
	}

	exist, err := c.FileIsExist(fpath, c.L)
	if err != nil {
		return c.Cache.SetErr(fpath, originalSVG, thumb_size, c.L.Typ.Error(logCP, logC02, err))
	}
	if !exist {
		return c.Cache.SetErr(fpath, originalSVG, thumb_size, c.L.Typ.Error(logCP, logC03, err))
	}

	if !originalSVG {

		source_width, source_height, err := c.Resolution.Get()
		if err != nil {
			return c.Cache.SetErr(fpath, originalSVG, thumb_size, c.L.Typ.Error(logCP, logC04, err))
		}

		if source_height < thumb_size || source_width < thumb_size {
			return c.Cache.SetErr(
				fpath, originalSVG, thumb_size, c.L.Typ.Error(logCP, 
					fmt.Sprintf(`Source: %vx%v, Preview: %v; %s`, source_width, source_height, thumb_size, logC05),
					err),
				)
		}

		if 	source_height < config_.ImageSourceResolutionMin || 
			source_height > config_.ImageSourceResolutionMax ||
			source_width < config_.ImageSourceResolutionMin || 
			source_width > config_.ImageSourceResolutionMax {

				return c.Cache.SetErr(
					fpath, originalSVG, thumb_size, c.L.Typ.Error(logCP, 
						fmt.Sprintf(`Min Resolution: %v, Max Resolution: %v, Current Resolution: %vx%v`,
						config_.ImageSourceResolutionMin,
						config_.ImageSourceResolutionMax,
						source_width,
						source_height,
					), logC06, err))
		}
	}

	c.Cache.SetErr(fpath, originalSVG, thumb_size, nil)
	return nil
}
