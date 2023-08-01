package checks_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 25 July 2023
 */
import (
	"errors"
	"testing"

	config_ "github.com/jhekau/favicon/internal/config"
	types_ "github.com/jhekau/favicon/internal/core/types"
	checks_ "github.com/jhekau/favicon/internal/service/convert/checks"
)


func TestCheckPreviewUnit( t *testing.T ) {

	for _, ts := range []struct{
		size int
		typ types_.FileType
		err error
	}{
		{ 0, 		types_.PNG(), errors.New(`error`) },
		{ 15, 		types_.ICO(), errors.New(`error`) },
		{ 16, 		types_.ICO(), nil },
		{ 10001, 	types_.PNG(), errors.New(`error`) },
		{ 10001,	types_.SVG(), nil },
		{ config_.ImagePreviewResolutionMin, 	types_.PNG(), nil },
		{ config_.ImagePreviewResolutionMin-1, 	types_.PNG(), errors.New(`error`) },
		{ config_.ImagePreviewResolutionMax, 	types_.PNG(), nil },
		{ config_.ImagePreviewResolutionMax+1, 	types_.PNG(), errors.New(`error`) },
	}{
		err := checks_.Preview{}.Check( ts.typ, ts.size)
		if (err == nil && ts.err != nil) || (err != nil && ts.err == nil) {
			t.Fatalf(`TestCheckPreviewUnit - status: data: %#v`, ts)
		}
	}

}

func TestCheckSourceCacheNotExistUnit( t *testing.T ) {

	check := checks_.CacheStatus{}

	data_not_exist := map[types_.FilePath]struct{
		err error
		exist bool
	}{
		`TestCheckSourceCacheNotExistUnit/1.jpg`: {nil, true},
		`TestCheckSourceCacheNotExistUnit/2.jpg`: {errors.New(`2.jpg`), true},
		`TestCheckSourceCacheNotExistUnit/3.jpg`: {errors.New(`3.jpg`), true},
	}
	
	// read
	for filepath, ts := range data_not_exist {
		status, err := check.Status(filepath, types_.ICO(), 16)
		if err != nil {
			t.Fatalf(`TestCheckSourceCacheNotExistUnit - error: data: %#v`, ts)
		}
		if status {
			t.Fatalf(`TestCheckSourceCacheNotExistUnit - status: data: %#v`, ts)
		}
	}

}

func TestCheckSourceCacheExistUnit( t *testing.T ) {

	check := checks_.CacheStatus{}

	data_not_exist := []struct{
		filepath 	types_.FilePath
		typ 		types_.FileType
		thumb_size 	int
		err 		error
		exist 		bool
	}{
		{`TestCheckSourceCacheExistUnit/1.jpg`, types_.PNG(), 16, nil, true},
		{`TestCheckSourceCacheExistUnit/1.jpg`, types_.ICO(), 16, errors.New(`1.jpg, ico, 16`), true},
		{`TestCheckSourceCacheExistUnit/1.jpg`, types_.ICO(), 20, errors.New(`1.jpg, ico, 20`), true},
		{`TestCheckSourceCacheExistUnit/2.jpg`, types_.PNG(), 16, errors.New(`2.jpg`), true},
		{`TestCheckSourceCacheExistUnit/3.jpg`, types_.PNG(), 16, errors.New(`3.jpg`), true},
	}
	
	// store
	for _, ts := range data_not_exist {
		check.SetErr(ts.filepath, ts.typ, ts.thumb_size, ts.err)
	}

	// read
	for _, ts := range data_not_exist {
		status, err := check.Status(ts.filepath, ts.typ, ts.thumb_size)
		if err != ts.err {
			t.Fatalf(`TestCheckSourceCacheExistUnit - error: data: %#v`, ts)
		}
		if !status {
			t.Fatalf(`TestCheckSourceCacheExistUnit - status: data: %#v`, ts)
		}
	}
}

type CheckSourceUnitCacheDisable struct{}
func (c CheckSourceUnitCacheDisable) Status(_ types_.FilePath, _ types_.FileType, _ int) (bool, error) {
	return false, nil
}
func (c CheckSourceUnitCacheDisable) SetErr(_ types_.FilePath, _ types_.FileType, _ int, err error) error {
	return err
}

type resolution struct{
	f func() (w int, h int, err error)
}
func (r resolution) Get() (w int, h int, err error){
	return r.f()
}

func TestCheckSourceUnit( t *testing.T ) {

	// enable and disable cache ******************
	type cache interface {
		Status(_ types_.FilePath, _ types_.FileType, _ int) (bool, error)
		SetErr(_ types_.FilePath, _ types_.FileType, _ int, _ error) error
	}

	cache_enable := checks_.CacheStatus{}
	cache_disable := CheckSourceUnitCacheDisable{}

	// test *************************************************************

	// init
	file_is_exist := struct{
		exist, not_exist, err func(fpath types_.FilePath) (bool, error)
	}{
		exist: 		func(fpath types_.FilePath) (bool, error){ return true, nil },
		not_exist: 	func(fpath types_.FilePath) (bool, error){ return false, nil },
		err: 		func(fpath types_.FilePath) (bool, error){ return false, errors.New(`error`) },
	}

	// testing
	for _, dt := range []struct{
		filepath         types_.FilePath
		typ 			 types_.FileType
		thumb_size 		 int
		file_is_exist 	 func(fpath types_.FilePath) (bool, error)
		file_resolution  func() (w int, h int, err error)
		cache 			 cache
		status_error 	 error
	}{
		{	`TestCheckSourceUnit/1.jpg`, // ошибка - нулевой размер ---------------------------------------
			types_.PNG(), 0, file_is_exist.exist, 
			func() (w int, h int, err error){ return 1, 1, nil },
			cache_disable,
			errors.New(`error`),
		},
		{	`TestCheckSourceUnit/2.jpg`, // ошибка - размер оригинала меньше, чем нарезаемая превьха ------
			types_.PNG(), 16, file_is_exist.exist, 
			func() (w int, h int, err error){ return 1, 1, nil },
			&cache_enable,
			errors.New(`error`),
		},
		{	`TestCheckSourceUnit/2.jpg`, // проверка работы кеша по предыдущему условию -------------------
			types_.PNG(), 16, file_is_exist.exist, 
			func() (w int, h int, err error){ return 1, 1, nil },
			&cache_enable,
			errors.New(`error`),
		},
		{	`TestCheckSourceUnit/3.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.exist, 
			func() (w int, h int, err error){ return 16, 16, nil },
			&cache_enable,
			nil,
		},
		{	`TestCheckSourceUnit/3.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.exist, 
			func() (w int, h int, err error){ return 16, 16, nil },
			&cache_enable,
			nil,
		},
		{	`TestCheckSourceUnit/4.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.not_exist, 
			func() (w int, h int, err error){ return 16, 16, nil },
			cache_disable,
			errors.New(`error`),
		},
		{	`TestCheckSourceUnit/5.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.err, 
			func() (w int, h int, err error){ return 16, 16, nil },
			cache_disable,
			errors.New(`error`),
		},
		{	`TestCheckSourceUnit/6.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.exist, 
			func() (w int, h int, err error){ return 16, 16, errors.New(`error`) },
			cache_disable,
			errors.New(`error`),
		},
	}{

		err := (&checks_.Source{
			Cache: dt.cache,
			FileIsExist: dt.file_is_exist,
			Resolution: resolution{dt.file_resolution},
		}).
		Check(dt.filepath, dt.typ, dt.thumb_size)

		if (err == nil && dt.status_error != nil) || (err != nil && dt.status_error == nil) {
			t.Fatalf(`TestCheckSourceUnit - error: filepath: %v err: '%v' data: %#v`, dt.filepath, err, dt)
		}
	}
}