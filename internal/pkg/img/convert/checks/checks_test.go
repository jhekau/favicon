package checks_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 25 July 2023
 */
import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	config_ "github.com/jhekau/favicon/internal/pkg/img/config"
	mock_logger_ "github.com/jhekau/favicon/internal/test/mocks/interfaces/logger"
	checks_ "github.com/jhekau/favicon/internal/pkg/img/convert/checks"
	storage_ "github.com/jhekau/favicon/interfaces/storage"
	types_ "github.com/jhekau/favicon/domain/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)


func TestCheckPreviewUnit( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
		logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
		 

		err := checks_.Preview{
			logs,
		}.Check( ts.typ, ts.size)

		if (err == nil && ts.err != nil) || (err != nil && ts.err == nil) {
			t.Fatalf(`TestCheckPreviewUnit - status: data: %#v`, ts)
		}
	}

}

func TestCheckSourceCacheNotExistUnit( t *testing.T ) {

	check := checks_.CacheStatus{}

	data_not_exist := map[storage_.StorageKey]struct{
		err error
		exist bool
	}{
		`TestCheckSourceCacheNotExistUnit/1.jpg`: {nil, true},
		`TestCheckSourceCacheNotExistUnit/2.jpg`: {errors.New(`2.jpg`), true},
		`TestCheckSourceCacheNotExistUnit/3.jpg`: {errors.New(`3.jpg`), true},
	}
	
	// read
	for filepath, ts := range data_not_exist {
		status, err := check.Status(filepath, false, 16)

		require.ErrorIs(t, err, nil, 
			fmt.Sprintf(`TestCheckSourceCacheNotExistUnit - error: data: %#v`, ts))

		require.False(t, status, 
			fmt.Sprintf(`TestCheckSourceCacheNotExistUnit - status: data: %#v`, ts) )
	}

}

func TestCheckSourceCacheExistUnit( t *testing.T ) {

	check := checks_.CacheStatus{}

	d := []struct{
		originalKey storage_.StorageKey
		originalSVG	bool
		thumb_size 	int
		err 		error
		exist 		bool
	}{
		{`TestCheckSourceCacheExistUnit/1.jpg`, false, 16, nil, true},
		{`TestCheckSourceCacheExistUnit/1.jpg`, false, 18, errors.New(`1.jpg, ico, 18`), true},
		{`TestCheckSourceCacheExistUnit/1.jpg`, false, 20, errors.New(`1.jpg, ico, 20`), true},
		{`TestCheckSourceCacheExistUnit/2.jpg`, false, 16, errors.New(`2.jpg`), true},
		{`TestCheckSourceCacheExistUnit/3.jpg`, false, 16, errors.New(`3.jpg`), true},
	}
	
	// store
	for _, ts := range d {
		check.SetErr(ts.originalKey, ts.originalSVG, ts.thumb_size, ts.err)
	}

	// read
	for _, ts := range d {
		status, err := check.Status(ts.originalKey, ts.originalSVG, ts.thumb_size)

		require.Equal(t, err, ts.err, 
			fmt.Sprintf(`TestCheckSourceCacheExistUnit - error: %v, data: %#v`, err, ts))

		require.True(t, status, 
			fmt.Sprintf(`TestCheckSourceCacheExistUnit - status: data: %#v`, ts) )
	}
}

type CheckSourceUnitCacheDisable struct{}
func (c CheckSourceUnitCacheDisable) Status(_ storage_.StorageKey, _ bool, _ int) (bool, error) {
	return false, nil
}
func (c CheckSourceUnitCacheDisable) SetErr(_ storage_.StorageKey, _ bool, _ int, err error) error {
	return err
}


type storageOBJ struct {
	key storage_.StorageKey
	isExists bool
	isExists_err error
}
func (s *storageOBJ) Key() storage_.StorageKey {
	return storage_.StorageKey(s.key)
}
func (s *storageOBJ) Reader() (io.ReadSeekCloser, error) {
	return nil, nil
}
func (s *storageOBJ) Writer() (io.WriteCloser, error) {
	return nil, nil
}
func (s *storageOBJ) IsExists() ( bool, error ) {
	return s.isExists, s.isExists_err
}
func (s *storageOBJ) ModTime() time.Time {
	return time.Time{}
}

type resolution struct{
	f func() (w int, h int, err error)
}
func (r resolution) Get(_ storage_.StorageOBJ) (w int, h int, err error){
	return r.f()
}

func TestCheckSourceUnit( t *testing.T ) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// enable and disable cache ******************
	type cache interface {
		Status(_ storage_.StorageKey, _ bool, _ int) (bool, error)
		SetErr(_ storage_.StorageKey, _ bool, _ int, _ error) error
	}

	cache_enable := checks_.CacheStatus{}
	cache_disable := CheckSourceUnitCacheDisable{}

	// test *************************************************************

	for _, dt := range []struct{
		storageOBJ       *storageOBJ
		originalSVG  	 bool
		thumb_size 		 int
		file_resolution  func() (w int, h int, err error)
		cache 			 cache
		status_error 	 error
	}{
		{	// ошибка - нулевой размер ---------------------------------------
			&storageOBJ{
				key:`TestCheckSourceUnit/1.jpg`,
				isExists: true,
				isExists_err: nil,
			},
			false, 0,
			func() (w int, h int, err error){ return 1, 1, nil },
			cache_disable,
			errors.New(`error`),
		},
		{	// ошибка - размер оригинала меньше, чем нарезаемая превьха ------
			&storageOBJ{
				key:`TestCheckSourceUnit/2.jpg`,
				isExists: true,
				isExists_err: nil,
			}, 
			false, 16, 
			func() (w int, h int, err error){ return 1, 1, nil },
			&cache_enable,
			errors.New(`error`),
		},
		{	// проверка работы кеша по предыдущему условию -------------------
			&storageOBJ{
				key:`TestCheckSourceUnit/2.jpg`,
				isExists: true,
				isExists_err: nil,
			}, 
			false, 16,
			func() (w int, h int, err error){ return 1, 1, nil },
			&cache_enable,
			errors.New(`error`),
		},
		{	// ------------------------------------------------------
			&storageOBJ{
				key:`TestCheckSourceUnit/3.jpg`,
				isExists: true,
				isExists_err: nil,
			},
			false, 16,
			func() (w int, h int, err error){ return 16, 16, nil },
			&cache_enable,
			nil,
		},
		{	// ------------------------------------------------------
			&storageOBJ{
				key:`TestCheckSourceUnit/3.jpg`,
				isExists: true,
				isExists_err: nil,
			}, 
			false, 16, 
			func() (w int, h int, err error){ return 16, 16, nil },
			&cache_enable,
			nil,
		},
		{	// ------------------------------------------------------
			&storageOBJ{
				key:`TestCheckSourceUnit/4.jpg`,
				isExists: false,
				isExists_err: nil,
			}, 
			false, 16, 
			func() (w int, h int, err error){ return 16, 16, nil },
			cache_disable,
			errors.New(`error`),
		},
		{	// ------------------------------------------------------
			&storageOBJ{
				key:`TestCheckSourceUnit/5.jpg`,
				isExists: false,
				isExists_err: errors.New(`error`),
			}, 
			false, 16, 
			func() (w int, h int, err error){ return 16, 16, nil },
			cache_disable,
			errors.New(`error`),
		},
		{	// ------------------------------------------------------
			&storageOBJ{
				key:`TestCheckSourceUnit/6.jpg`,
				isExists: true,
				isExists_err: nil,
			}, 
			false, 16, 
			func() (w int, h int, err error){ return 16, 16, errors.New(`error`) },
			cache_disable,
			errors.New(`error`),
		},
	}{
		logs := mock_logger_.NewMockLogger(ctrl)
	logs.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
		 

		err := (&checks_.Source{
			L: logs,
			Cache: dt.cache,
			Resolution: resolution{ dt.file_resolution },
		}).
		Check(dt.storageOBJ, dt.originalSVG, dt.thumb_size)

		if (err == nil && dt.status_error != nil) || (err != nil && dt.status_error == nil) {
			t.Fatalf(`TestCheckSourceUnit - error: filepath: %v err: '%v' data: %#v`, dt.storageOBJ.Key(), err, dt)
		}
	}
}