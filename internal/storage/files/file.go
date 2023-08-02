package files

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	"os"

	// err_ "github.com/jhekau/favicon/internal/core/err"
	logger_ "github.com/jhekau/favicon/internal/core/logger"
	types_ "github.com/jhekau/favicon/internal/core/types"
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

type File struct{
	L *logger_.Logger
	fpath types_.FilePath
	f *os.File
}

func (s *File) Read() (*os.File, error) {
	f, err := osOpen(s.fpath.String())
	if err != nil {
		return nil, s.L.Typ.Error(logP, logS02, err)
	}
	return f, nil
}
func (s *File) Close() error {
	if s.f == nil {
		return nil
	}
	err := s.f.Close()
	if err != nil {
		return s.L.Typ.Error(logP, logS01, err)
	}
	return nil
}

func (s *File) IsExists() ( bool, error ) {

	if f, err := osStat(s.fpath.String()); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, s.L.Typ.Error(logP, logS03, err)

	} else if f.IsDir() {
		return false, s.L.Typ.Error(logP, logS04)
	}

	return true, nil
}

// TODO ?
// func (s *File) Write()....
// func (s *File) Delete()....


func New( l *logger_.Logger, f types_.FilePath ) *File {
	return &File{
		L: l,
		fpath: f,
	}
}

