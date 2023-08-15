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

	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"
)

const (
	logP = `/internal/storage/files/stat.go`
	logS01 = `S01: can't get path`
	logS02 = `S02: failed open file`
	logS03 = `S03: os stat source image`
	logS04 = `S04: save thumb image is a folder`
	logS05 = `S05: failed open file`
	logS06 = `S06: os stat source image`
	logS07 = `S07: can't get path`
)

// for test 
var (
	osOpen = os.Open
	osStat = os.Stat
)

var folderIcons = `icons` // default

func SetFolderIcons(f string) {
	folderIcons = f
}

// storage object
type file struct{
	l logger_.Logger
	key string
}

// storage
type Files struct{
	L logger_.Logger
}

func (s *file) Reader() (io.ReadSeekCloser, error) {

	filename, err := filepath.Abs(filepath.Join(folderIcons, s.key))
	if err != nil {
		return nil, s.l.Error(logP, logS01, err)
	}
	f, err := osOpen(filename)
	if err != nil {
		return nil, s.l.Error(logP, logS02, err)
	}
	return f, nil
}

func (s *file) Writer() (io.WriteCloser, error){
	
	filename, err := filepath.Abs(filepath.Join(folderIcons, s.key))
	if err != nil {
		return nil, s.l.Error(logP, logS07, err)
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return nil, s.l.Error(logP, logS05, err)
	}
	return f, nil
}

func (s *file) IsExists() ( bool, error ) {

	if s == nil {
		return false, nil
	}
	
	if f, err := osStat( filepath.Join(folderIcons, s.key) ); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, s.l.Error(logP, logS03, err)

	} else if f.IsDir() {
		return false, s.l.Error(logP, logS04)
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
	st, err := osStat(filepath.Join(folderIcons, s.key))
	if err != nil {
		s.l.Error(logP, logS06, err)
		return time.Time{}
	}
	
	return st.ModTime()
}

// TODO ?
// func (s *File) Delete()....

// получаем интерфейсы на объект в storage
func (fl Files) NewObject( key any ) (storage_.StorageOBJ, error) {
	return &file{
		l: fl.L,
		key: key.(string),
	}, nil
}

