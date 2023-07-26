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
		status bool
	}{
		{ 0, 		types_.PNG(), false },
		{ 15, 		types_.ICO(), false },
		{ 16, 		types_.ICO(), true },
		{ 10001, 	types_.PNG(), false },
		{ 10001,	types_.SVG(), true },
		{ config_.ImagePreviewResolutionMin, 	types_.PNG(), true },
		{ config_.ImagePreviewResolutionMin-1, 	types_.PNG(), false },
		{ config_.ImagePreviewResolutionMax, 	types_.PNG(), true },
		{ config_.ImagePreviewResolutionMax+1, 	types_.PNG(), false },
	}{
		status, err := checks_.Preview{}.Check( ts.typ, ts.size)
		if err != nil {
			t.Fatalf(`CheckPreviewTest - error: data: %#v`, ts)
		}
		if status != ts.status {
			t.Fatalf(`CheckPreviewTest - status: data: %#v`, ts)
		}
	}

}

func TestCheckSourceCacheNotExistUnit( t *testing.T ) {

	check := checks_.ExportCheckCache()

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
			t.Fatalf(`TestCheckSourceCacheExist - error: data: %#v`, ts)
		}
		if status {
			t.Fatalf(`TestCheckSourceCacheExist - status: data: %#v`, ts)
		}
	}

}

func TestCheckSourceCacheExistUnit( t *testing.T ) {

	check := checks_.ExportCheckCache()

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
			t.Fatalf(`TestCheckSourceCacheExist - error: data: %#v`, ts)
		}
		if !status {
			t.Fatalf(`TestCheckSourceCacheExist - status: data: %#v`, ts)
		}
	}
}

func TestCheckSourceUnit( t *testing.T ) {

	// backup check cache function and disable cache ******************
	check_cache := func() (fn struct { defaulT func(); disable func() }) {

		check := checks_.ExportCheckCache()
		def := *check

		check_disable := struct{
			Status func(fpath types_.FilePath, source_typ types_.FileType, thumb_size int) (bool, error)
			SetErr func(fpath types_.FilePath, source_typ types_.FileType, thumb_size int, err error) error
		}{
			// disable cache
			Status: func(fpath types_.FilePath, source_typ types_.FileType, thumb_size int) (bool, error) {
				return false, nil
			},
			SetErr: func(fpath types_.FilePath, source_typ types_.FileType, thumb_size int, err error) error {
				return err
			},
		}

		fn.disable = func() {
			*check = check_disable
		}
		fn.defaulT = func() {
			*check = def
		}
		return
	}

	// restore check cache function and enable cache
	defer check_cache().defaulT()	

	// test *************************************************************

	// init
	file_is_exist := struct{
		exist, not_exist, err func(fpath types_.FilePath) (bool, error)
	}{
		exist: 		func(fpath types_.FilePath) (bool, error){ return true, nil },
		not_exist: 	func(fpath types_.FilePath) (bool, error){ return false, nil },
		err: 		func(fpath types_.FilePath) (bool, error){ return false, errors.New(`error`) },
	}

	source_check := checks_.Source{}

	// testing
	for _, dt := range []struct{
		filepath         types_.FilePath
		typ 			 types_.FileType
		thumb_size 		 int
		file_is_exist 	 func(fpath types_.FilePath) (bool, error)
		file_resolution  func(fpath types_.FilePath) (w int, h int, err error)
		check_cache_used func()
		status_error 	 error
	}{
		{	`TestCheckSourceUnit/1.jpg`, // ------------------------------------------------------
			types_.PNG(), 0, file_is_exist.exist, 
			func(fpath types_.FilePath) (w int, h int, err error){ return 1, 1, nil },
			check_cache().disable,
			errors.New(`error`),
		},
		{	`TestCheckSourceUnit/2.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.exist, 
			func(fpath types_.FilePath) (w int, h int, err error){ return 1, 1, nil },
			check_cache().defaulT,
			errors.New(`error`),
		},
		{	`TestCheckSourceUnit/2.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.exist, 
			func(fpath types_.FilePath) (w int, h int, err error){ return 1, 1, nil },
			check_cache().defaulT,
			errors.New(`error`),
		},
		{	`TestCheckSourceUnit/3.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.exist, 
			func(fpath types_.FilePath) (w int, h int, err error){ return 16, 16, nil },
			check_cache().defaulT,
			nil,
		},
		{	`TestCheckSourceUnit/3.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.exist, 
			func(fpath types_.FilePath) (w int, h int, err error){ return 16, 16, nil },
			check_cache().defaulT,
			nil,
		},
		{	`TestCheckSourceUnit/4.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.not_exist, 
			func(fpath types_.FilePath) (w int, h int, err error){ return 16, 16, nil },
			check_cache().disable,
			errors.New(`error`),
		},
		{	`TestCheckSourceUnit/5.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.err, 
			func(fpath types_.FilePath) (w int, h int, err error){ return 16, 16, nil },
			check_cache().disable,
			errors.New(`error`),
		},
		{	`TestCheckSourceUnit/6.jpg`, // ------------------------------------------------------
			types_.PNG(), 16, file_is_exist.exist, 
			func(fpath types_.FilePath) (w int, h int, err error){ return 16, 16, errors.New(`error`) },
			check_cache().disable,
			errors.New(`error`),
		},
	}{
		// swith data cache check used
		dt.check_cache_used()

		source_check.FileIsExist = dt.file_is_exist
		source_check.FileResolution = dt.file_resolution

		err := source_check.Check(dt.filepath, dt.typ, dt.thumb_size)
		if (err == nil && dt.status_error != nil) || (err != nil && dt.status_error == nil) {
			t.Fatalf(`TestCheckSourceUnit - error: filepath: %v err: '%v' data: %#v`, dt.filepath, err, dt)
		}
	}

	/*
	{ 0, 		types_.PNG(), false },
	{ 15, 		types_.ICO(), false },
	{ 16, 		types_.ICO(), true },
	{ 10001, 	types_.PNG(), false },
	{ 10001,	types_.SVG(), true },
	{ config_.ImagePreviewResolutionMin, 	types_.PNG(), true },
	{ config_.ImagePreviewResolutionMin-1, 	types_.PNG(), false },
	{ config_.ImagePreviewResolutionMax, 	types_.PNG(), true },
	{ config_.ImagePreviewResolutionMax+1, 	types_.PNG(), false },
	*/

}