package files

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	"io"
	"os"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	typ_ "github.com/jhekau/favicon/internal/core/types"
	storage_ "github.com/jhekau/favicon/pkg/models/storage"
)

const (
	logP = `/internal/storage/files/stat.go`
	logS01 = `S01: failed close file`
	logS02 = `S02: error open file`
	logS03 = `S03: os stat source image`
	logS04 = `S04: save thumb image is a folder`
	logS05 = `S05: failed open file`
)

// for test 
var (
	osOpen = os.Open
	osStat = os.Stat
)



// storage object
type file struct{
	l *logger_.Logger
	filepath typ_.FilePath
	f *os.File
}

// storage
type Files struct{
	L *logger_.Logger
}

func (s *file) Reader() (io.ReadCloser, error) {
	f, err := osOpen(s.filepath.String())
	if err != nil {
		return nil, s.l.Typ.Error(logP, logS02, err)
	}
	return f, nil
}
func (s *file) Writer() (io.WriteCloser, error){
	f, err := os.OpenFile(s.filepath.String(), os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return nil, s.l.Typ.Error(logP, logS05, err)
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

	if s == nil {
		return false, nil
	}
	
	if f, err := osStat(s.filepath.String()); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, s.l.Typ.Error(logP, logS03, err)

	} else if f.IsDir() {
		return false, s.l.Typ.Error(logP, logS04)
	}

	return true, nil
}

func (s *file) Key() storage_.StorageKey {
	return storage_.StorageKey(s.filepath.String())
}

// TODO ?
// func (s *File) Delete()....

// получаем интерфейсы на объект в storage
func (fl *Files) NewObject( filepath typ_.FilePath ) *file {
	return &file{
		l: fl.L,
		filepath: filepath,
	}
}

