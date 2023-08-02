package files_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 28 July 2023
 */
import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	logger_mock_ "github.com/jhekau/favicon/internal/core/logger/mock"
	files_ "github.com/jhekau/favicon/internal/storage/files"
)

type file_info struct {
	isDir bool
}
func(f *file_info) Name() string { return `` }
func(f *file_info) Size() int64 { return 0 }
func(f *file_info) Mode() os.FileMode { return 0 }
func(f *file_info) ModTime() time.Time { return time.Time{} }
func(f *file_info) IsDir() bool { return f.isDir }
func(f *file_info) Sys() any { return nil }

func TestIsExists(t *testing.T) {

	//
	backupOsStat := *files_.OsStat
	backupOsOpen := *files_.OsOpen
	defer func(){
		*files_.OsStat = backupOsStat
		*files_.OsOpen = backupOsOpen
	}()

	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}

	for _, d := range []struct{
		osOpen func(_ string) (*os.File, error)
		osStat func(_ string) (fs.FileInfo, error)
		resultIsExist bool
		resultError error
	}{
		{	// true, nil, !IsDir -> true, nil
			func(_ string) (*os.File, error) { return nil, nil },
			func(_ string) (fs.FileInfo, error) { return &file_info{}, nil }, // exist, not error
			true, nil,
		},
		{	// true, nil, IsDir -> false, error				!если директория
			func(_ string) (*os.File, error) { return nil, nil },
			func(_ string) (fs.FileInfo, error) { return &file_info{ true }, nil }, // exist, not error
			false, logger.Typ.Error(files_.LogP, files_.LogS04),
		},
		{	// false, error(os.ErrNotExist) -> false, nil	!если файла нет
			func(_ string) (*os.File, error) { return nil, nil },
			func(_ string) (fs.FileInfo, error) { return &file_info{}, os.ErrNotExist }, // exist, not error
			false, nil,
		},
		{	// false, error(error) -> false, error
			func(_ string) (*os.File, error) { return nil, nil },
			func(_ string) (fs.FileInfo, error) { return &file_info{}, errors.New(`error`) }, // exist, error
			false, logger.Typ.Error(files_.LogP, files_.LogS03, errors.New(`error`)),
		},
	}{
		*files_.OsStat = d.osStat
		*files_.OsOpen = d.osOpen

		isExist, err := files_.New(logger, ``).IsExists()

		require.Equal(t, err, d.resultError, fmt.Sprintf(
			`error: isExists: '%v', resultIsExist: '%v', err: '%v', data: %#v`, isExist, d.resultIsExist, err, d))

		require.Equal(t, isExist, d.resultIsExist, fmt.Sprintf(
			`status: isExists: '%v', resultIsExist: '%v', err: '%v', data: %#v`, isExist, d.resultIsExist, err, d)) 

	}
}

func TestRead(t *testing.T) {

	//
	backupOsOpen := *files_.OsOpen
	defer func(){
		*files_.OsOpen = backupOsOpen
	}()

	logger := &logger_.Logger{
		Typ: &logger_mock_.LoggerErrorf{},
	}

	for _, d := range []struct{
		osOpen func(_ string) (*os.File, error)
		resultIsExist *os.File
		resultError error
	}{
		{ 
			func(_ string) (*os.File, error) { return nil, nil }, 
			nil, 
			nil },
		{ 
			func(_ string) (*os.File, error) { return &os.File{}, nil }, 
			&os.File{}, 
			nil },
		{ 
			func(_ string) (*os.File, error) { return nil, errors.New(`error`) },
			nil, 
			logger.Typ.Error(files_.LogP, files_.LogS02, errors.New(`error`)) },
	}{
		*files_.OsOpen = d.osOpen

		f, err := files_.New(logger, ``).Read()

		require.Equal(t, err, d.resultError, 
			fmt.Sprintf(`error: f: '%v', err: '%v', data: %#v`, f, err, d)) 

		require.EqualValues(t, fmt.Sprint(f), fmt.Sprint(d.resultIsExist),
			fmt.Sprintf(`file status: f: '%v', err: '%v', data: %#v`, f, err, d))
			
	}

}

