package favicon

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 09 March 2023
 */
import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	types_ "github.com/jhekau/favicon/internal/core/types"
	convert_ "github.com/jhekau/favicon/internal/service/convert"
	checks_ "github.com/jhekau/favicon/internal/service/convert/checks"
	converters_ "github.com/jhekau/favicon/internal/service/convert/converters"
	defaults_ "github.com/jhekau/favicon/internal/service/defaults"
	converter_exec_anthonynsimon_ "github.com/jhekau/favicon/internal/service/img/converter/anthonynsimon"
	resolution_ "github.com/jhekau/favicon/internal/service/img/resolution"
	manifest_ "github.com/jhekau/favicon/internal/service/manifest"
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
	files_ "github.com/jhekau/favicon/internal/storage/files"
)

const (
	logTP  = `pkg/thumbs.go`
	logT01 = `T01: get serve file`
	logT02 = `T02: get file manifest`
	logT03 = `T03: get file thumb`
	logT04 = `T04: get file thumb`
	logT05 = `T05: get file manifest`
	logT06 = `T06: error read original`
)

var (

	TypePNG = types_.PNG()
	TypeICO = types_.ICO()
	// TypeSVG = types_.SVG()

	// ~~ interface ~~ 

	DefaultList = default_list

	/* -------------------------------------------------------------------------------------------
	 * CustomList(
	 *	ThumbNEW().SetHREF(`touch-icon-iphone.png`).SetSize(180).SetTypeImage(favicon.TypePNG),
	 *	ThumbNEW().SetHREF(`/icon-192.png`).SetSize(192).SetTypeImage(favicon.TypePNG),
	 *	ThumbNEW().SetHREF(`/icon-512.png`).SetSize(512).SetTypeImage(favicon.TypePNG),
	 * )
	 * ------------------------------------------------------------------------------------------- */
	CustomList = custom_list

	ThumbNEW = func() *thumb_.Thumb {
		return &thumb_.Thumb{}
	}
)

type Converter interface{
	Do(source, save types_.FilePath, originalSVG bool, typ types_.FileType, size_px int) error
}

type StorageObject interface{
	Read() (*os.File, error)
	Close() error
	IsExists() ( bool, error )
}

type Thumbs struct {
	s sync.RWMutex
	l interface { // logger
        Info(path string, messages ...interface{}) error
        Alert(path string, messages ...interface{}) error
        Error(path string, messages ...interface{}) error
    }
	storage interface {
		Object( key func() string ) StorageObject
	}
	source_svg types_.FilePath
	source_img types_.FilePath
	folder_work types_.Folder
	thumbs map[types_.URLHref/*clear*/]*thumb_.Thumb
	manifest manifest_.Manifest
}

func (t *Thumbs) Append(thumb *thumb_.Thumb) *Thumbs {
	return t.append(thumb)
}
func (t *Thumbs) SetFolderWork( folder string ) *Thumbs {
	return t.set_folder_work(folder)
}
func (t *Thumbs) SetFilepathSourceSVG( fpath string ) *Thumbs {
	return t.set_filepath_source_svg( fpath )
}
func (t *Thumbs) SetFilepathSourceIMG( fpath string ) *Thumbs {
	return t.set_filepath_source_img( fpath )
}
func (t *Thumbs) Handle() {
	t.handle()
}


// использует конвертер и систему хранения изображений по умолчанию
func (t *Thumbs) ServeFile( url_ *url.URL ) ( fpath string, exists bool, err error ) {

	converter := converter_exec_anthonynsimon_.Exec{}
	tb, exist := t.getThumb(types_.URLHref(url_.Path))
	if !exist {
		return ``, false, nil
	}

	// storage := files_.Files{ t.l }

	storageObj := t.storage.Object( tb.GetOriginalKey )

	return t.serve_file(url_, &convert_.Converter{
		Converters: []convert_.ConverterT{
			&converters_.ConverterPNG{ConverterExec: &converter},
			&converters_.ConverterICO{ConverterExec: &converter},
		},
		CheckPreview: checks_.Preview{},
		CheckSource: &checks_.Source{
			Cache: &checks_.CacheStatus{},
			StorageObj: storageObj,
			Resolution: &resolution_.Resolution{
				L: t.l,
			} ,
		},
	})
}
func (t *Thumbs) TagsHTML() string {
	return t.tags_html()
}


func (t *Thumbs) append(thumb *thumb_.Thumb) *Thumbs {
	t.s.Lock()
	defer t.s.Unlock()
	
	t.thumbs[thumb.GetHREFClear()] = thumb
	return t
}

func (t *Thumbs) set_folder_work( folder string ) *Thumbs {
	t.s.Lock()
	defer t.s.Unlock()

	t.folder_work = types_.Folder(folder)
	return t
}
func (t *Thumbs) get_folder_work() types_.Folder {
	t.s.RLock()
	defer t.s.RUnlock()

	return t.folder_work
}
func (t *Thumbs) set_filepath_source_svg( fpath string ) *Thumbs {
	t.s.Lock()
	defer t.s.Unlock()

	t.source_svg = types_.FilePath(fpath)
	return t
}
func (t *Thumbs) get_filepath_source_svg() types_.FilePath {
	t.s.RLock()
	defer t.s.RUnlock()

	return t.source_svg
}
func (t *Thumbs) set_filepath_source_img( fpath string ) *Thumbs {
	t.s.Lock()
	defer t.s.Unlock()

	t.source_img = types_.FilePath(fpath)
	return t
}
func (t *Thumbs) get_filepath_source_img() types_.FilePath {
	t.s.RLock()
	defer t.s.RUnlock()

	return t.source_img
}

func (t *Thumbs) handle() {
	http.Handle(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		switch r.Method {
		case http.MethodGet, http.MethodPost:
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		fpath, exists, err := t.ServeFile(r.URL)
		if err != nil {
			log.Println(t.l.Typ.Error(logT01, err))
			w.WriteHeader(http.StatusInternalServerError)
		} else if !exists {
			w.WriteHeader(http.StatusNotFound)
		} else {
			http.ServeFile(w, r, fpath)
		}
	}))
}

func (t *Thumbs) serve_file( url_ *url.URL, conv Converter ) ( fpath string, exists bool, err error ) {

	if manifest, exists, err := t.server_file_manifest(url_); err != nil {
		return ``, false, t.l.Typ.Error(logTP, logT02, err)
	} else if exists {
		return manifest.String(), true, nil
	}

	if thumb, exists, err := t.server_file_thumb(url_, conv); err != nil {
		return ``, false, t.l.Typ.Error(logTP, logT03, err)
	} else if exists {
		return thumb.String(), true, nil
	}
	return ``, false, nil
}

func (t *Thumbs) server_file_thumb( url_ *url.URL, conv Converter ) (fpath types_.FilePath, exists bool, err error) {

	t.s.RLock()
	thumb, exists := thumb_.URLExists(url_, t.thumbs)
	t.s.RUnlock()

	if !exists {
		return ``, false, nil
	}

	fpath, err = thumb.GetFile(
		t.get_folder_work(),
		t.get_filepath_source_img(),
		t.get_filepath_source_svg(),
		conv,
	)
	if err != nil {
		return ``, false, t.l.Typ.Error(logTP, logT04, err)
	}

	return fpath, true, nil
}

func (t *Thumbs) server_file_manifest( url_ *url.URL ) (manifest types_.FilePath, exists bool, err error) {

	if !t.manifest.URLExists(url_) {
		return ``, false, nil
	}

	manifest, exists, err = t.manifest.GetFile(t.get_folder_work(), t.thumbs)
	if err != nil {
		return ``, false, t.l.Typ.Error(logTP, logT05, err)
	}
	return manifest, exists, nil
}

// func (T *Thumbs) manifest_url_exists( URLpath string ) ( fpath string, exists bool, err error )
// func (t *Thumbs) get_thumbs() map[types_.URLHref/*clear*/]*thumb_.Thumb

func (t *Thumbs) tags_html() string {

	var tags strings.Builder

	t.s.RLock()
	for _, thumb := range t.thumbs {
		tags.WriteString(thumb.GetTAG())
	}
	tags.WriteString(t.manifest.GetTAG())
	t.s.RUnlock()

	return tags.String()
}

func (t *Thumbs) getThumb(u types_.URLHref) (*thumb_.Thumb, bool) {
	tb, ok := t.thumbs[u]
	return tb, ok
}

func default_list() *Thumbs {
	return &Thumbs{
		thumbs: func() map[types_.URLHref]*thumb_.Thumb {
			m := map[types_.URLHref]*thumb_.Thumb{}
			for _, thumb := range defaults_.Defaults() {
				m[thumb.GetHREFClear()] = thumb
			}
			return m
		}(),
		manifest: manifest_.Default(),
	}
}

func custom_list( thumbs ...*thumb_.Thumb ) *Thumbs {
	return &Thumbs{
		thumbs: func() map[types_.URLHref]*thumb_.Thumb {
			m := map[types_.URLHref]*thumb_.Thumb{}
			for _, thumb := range thumbs {
				m[thumb.GetHREFClear()] = thumb
			}
			return m
		}(),
	}
}
