package files_test

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 28 July 2023
 */
import (
	"errors"
	"io/fs"
	"os"
	"testing"
	"time"

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
	backup := *files_.OsStat
	defer func(){
		*files_.OsStat = backup
	}()

	for _, d := range []struct{
		osStat func(_ string) (fs.FileInfo, error)
		resultIsExist bool
		resultError error
	}{
		{	// true, nil, !IsDir -> true, nil
			func(_ string) (fs.FileInfo, error) { return &file_info{}, nil }, // exist, not error
			true, nil,
		},
		{	// true, nil, IsDir -> false, error				!если директория
			func(_ string) (fs.FileInfo, error) { return &file_info{ true }, nil }, // exist, not error
			false, errors.New(`error`),
		},
		{	// false, error(os.ErrNotExist) -> false, nil	!если файла нет
			func(_ string) (fs.FileInfo, error) { return &file_info{}, os.ErrNotExist }, // exist, not error
			false, nil,
		},
		{	// false, error(error) -> false, error
			func(_ string) (fs.FileInfo, error) { return &file_info{}, errors.New(`error`) }, // exist, error
			false, errors.New(`error`),
		},
	}{
		*files_.OsStat = d.osStat

		isExist, err := files_.IsExists(``)
		if (err == nil && d.resultError != nil) || (err != nil && d.resultError == nil) {
			t.Fatalf(`TestIsExists - error: '%v' data: %#v`, err, d)
		}
		if isExist != d.resultIsExist {
			t.Fatalf(`TestIsExists - isExists: '%v' data: %#v`, err, d)
		}
	}


}