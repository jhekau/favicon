package files

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	"os"

	// err_ "github.com/jhekau/favicon/internal/core/err"
	logger_ "github.com/jhekau/favicon/internal/core/logger"
	// types_ "github.com/jhekau/favicon/internal/core/types"
	domain_ "github.com/jhekau/favicon/pkg/domain"
)

const (
	logP = `/internal/storage/files/stat.go`
	logS01 = `S01: failed close file`
	logS02 = `S02: error open file`
	logS03 = `S03: os stat source image`
	logS04 = `S04: save thumb image is a folder`
)

// for test 
var (
	osOpen = os.Open
	osStat = os.Stat
)

// storage object
type file struct{
	l *logger_.Logger
	filepath func() string
	f *os.File
}

// storage
type Files struct{
	L *logger_.Logger
}

func (s *file) Read() (*os.File, error) {
	f, err := osOpen(s.filepath())
	if err != nil {
		return nil, s.l.Typ.Error(logP, logS02, err)
	}
	return f, nil
}
func (s *file) Close() error {
	if s.f == nil {
		return nil
	}
	err := s.f.Close()
	if err != nil {
		return s.l.Typ.Error(logP, logS01, err)
	}
	return nil
}

func (s *file) IsExists() ( bool, error ) {

	if f, err := osStat(s.filepath()); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, s.l.Typ.Error(logP, logS03, err)

	} else if f.IsDir() {
		return false, s.l.Typ.Error(logP, logS04)
	}

	return true, nil
}

func (s *file) Key() domain_.StorageKey {
	return domain_.StorageKey(s.filepath())
}

// TODO ?
// func (s *File) Write()....
// func (s *File) Delete()....

// получаем интерфейсы на объект в storage
func (fl *Files) Object( key func() string ) *file {
	return &file{
		l: fl.L,
		filepath: key,
	}
}

