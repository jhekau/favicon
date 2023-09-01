package files

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	"io"
	"os"
	"path/filepath"
	"time"

	err_ "github.com/jhekau/favicon/internal/pkg/err"
	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"
)

const (
	logP = `/internal/storage/files/file.go`
	logS01 = `S01: can't get path`
	logS02 = `S02: failed open file`
	logS03 = `S03: os stat source image`
	logS04 = `S04: save thumb image is a folder`
	logS05 = `S05: failed open file`
	logS06 = `S06: os stat source image`
	logS07 = `S07: can't get path`
	logS08 = `S08: can't create dir`
)

// for test 
var (
	osOpen = os.Open
	osStat = os.Stat
)

const dirIconsDefault = `icons` // default

//
func NewStorage(dirIcons string, logger logger_.Logger) *storage {
	return &storage{
		dir: dirIcons,
		l: logger,
	}
}

// storage object
type file struct{
	l logger_.Logger
	key string
	dir string
}

func (s *file) Reader() (io.ReadSeekCloser, error) {

	fpath, err := filepath.Abs(filepath.Join(s.dir, s.key))
	if err != nil {
		return nil, err_.Err(s.l, logP, logS01, err)
	}
	f, err := osOpen(fpath)
	if err != nil {
		return nil, err_.Err(s.l, logP, logS02, err)
	}
	return f, nil
}

func (s *file) Writer() (io.WriteCloser, error){
	
	fpath, err := filepath.Abs(filepath.Join(s.dir, s.key))
	if err != nil {
		return nil, err_.Err(s.l, logP, logS07, err)
	}
	
	err = os.MkdirAll(filepath.Dir(fpath), 0775)
	if err != nil {
		return nil, err_.Err(s.l, logP, logS08, err)
	}
	
	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err_.Err(s.l, logP, logS05, err)
	}
	return f, nil
}

func (s *file) IsExists() ( bool, error ) {

	if s == nil {
		return false, nil
	}
	
	if f, err := osStat( filepath.Join(s.dir, s.key) ); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err_.Err(s.l, logP, logS03, err)

	} else if f.IsDir() {
		return false, err_.Err(s.l, logP, logS04)
	}

	return true, nil
}

func (s *file) Key() storage_.StorageKey {
	return storage_.StorageKey(s.key)
}

func (s *file) ModTime() time.Time {

	if s.key == `` {
		return time.Time{}
	}
	st, err := osStat(filepath.Join(s.dir, s.key))
	if err != nil {
		err_.Err(s.l, logP, logS06, err)
		return time.Time{}
	}
	
	return st.ModTime()
}

// TODO ?
// func (s *File) Delete()....


// storage
type storage struct{
	dir string
	l logger_.Logger
}

// получаем интерфейсы на объект в storage
func (s *storage) NewObject( key any ) (storage_.StorageOBJ, error) {
	return &file{
		l: s.l,
		key: key.(string),
		dir: s.dir,
	}, nil
}

func (s *storage) SetDirDefault() *storage {
	s.dir = dirIconsDefault
	return s
}

