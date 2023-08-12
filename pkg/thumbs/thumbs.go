package favicon

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 09 March 2023
 */
import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	logger_default_ "github.com/jhekau/favicon/internal/core/logs/default"
	typ_ "github.com/jhekau/favicon/internal/core/types"
	convert_ "github.com/jhekau/favicon/internal/service/convert"
	checks_ "github.com/jhekau/favicon/internal/service/convert/checks"
	converters_ "github.com/jhekau/favicon/internal/service/convert/converters"
	defaults_ "github.com/jhekau/favicon/internal/service/defaults"
	converter_exec_anthonynsimon_ "github.com/jhekau/favicon/internal/service/img/converter/anthonynsimon"
	resolution_ "github.com/jhekau/favicon/internal/service/img/resolution"
	manifest_ "github.com/jhekau/favicon/internal/service/manifest"
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
	files_ "github.com/jhekau/favicon/internal/storage/files"

	types_ "github.com/jhekau/favicon/pkg/core/types"
	converter_ "github.com/jhekau/favicon/pkg/models/converter"
	logger_ "github.com/jhekau/favicon/pkg/models/logger"
	storage_ "github.com/jhekau/favicon/pkg/models/storage"
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

// создание пустого набора превьюх для одного оригинального изображения
func NewThumbs() *Thumbs {
	logger := &logger_default_.Logger{}
	return &Thumbs{
		l: logger,
		storage: files_.Files{L: logger},
		conv: &convert_.Converter{
			L: logger,
			Converters: []converter_.ConverterTyp{
				&converters_.ConverterICO{
					ConverterExec: &converter_exec_anthonynsimon_.Exec{},
				},
				&converters_.ConverterPNG{
					ConverterExec: &converter_exec_anthonynsimon_.Exec{},
				},
			},
			CheckPreview: checks_.Preview{},
			CheckSource: &checks_.Source{
				L: logger,
				Cache: &checks_.CacheStatus{},
				Resolution: &resolution_.Resolution{
					L: logger,
				},
			},
		},
	}
}

// создание набора по умолчанию превьюх для оригинала
func NewThumbs_Defaults() *Thumbs {
	t := NewThumbs()
	return t
}

type Thumbs struct {
	mu          sync.RWMutex
	l           logger_.Logger
	storage     storage_.Storage
	conv 		converter_.Converter
	source_svg  typ_.FilePath
	source_img  typ_.FilePath
	folder_work typ_.Folder
	thumbs      map[typ_.URLHref]*thumb_.Thumb
	manifest    manifest_.Manifest
}

// возможность заменить логгер на собственную реализацию
func (t *Thumbs) LoggerSet( l logger_.Logger ) {
	t.mu.Lock()
	t.l = l
	t.mu.Unlock()
}

// возможность заменить систему хранения на собственную реализацию
func (t *Thumbs) StorageSet( s storage_.Storage ) {
	t.mu.Lock()
	t.storage = s
	t.mu.Unlock()
}

// возможность заменить конвертер на собственную реализацию
func (t *Thumbs) ConvertSet( conv converter_.Converter ) {
	t.mu.Lock()
	t.conv = conv
	t.mu.Unlock()
}



func (t *Thumbs) Append(thumb *thumb_.Thumb) *Thumbs {
	return t.append(thumb)
}
func (t *Thumbs) SetFolderWork(folder string) *Thumbs {
	return t.set_folder_work(folder)
}
func (t *Thumbs) SetFilepathSourceSVG(fpath string) *Thumbs {
	return t.set_filepath_source_svg(fpath)
}
func (t *Thumbs) SetFilepathSourceIMG(fpath string) *Thumbs {
	return t.set_filepath_source_img(fpath)
}
func (t *Thumbs) Handle() {
	t.handle()
}



// 
func (t *Thumbs) ServeFile(url_ *url.URL) (fpath string, exists bool, err error) {
	return t.serve_file(url_)
}



func (t *Thumbs) TagsHTML() string {
	return t.tags_html()
}

func (t *Thumbs) append(thumb *thumb_.Thumb) *Thumbs {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.thumbs[thumb.GetHREF()] = thumb
	return t
}

func (t *Thumbs) set_folder_work(folder string) *Thumbs {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.folder_work = typ_.Folder(folder)
	return t
}
func (t *Thumbs) get_folder_work() typ_.Folder {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.folder_work
}

/*
... это всё в storage
func (t *Thumbs) set_filepath_source_svg(fpath string) *Thumbs {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.source_svg = typ_.FilePath(fpath)
	return t
}
func (t *Thumbs) get_filepath_source_svg() typ_.FilePath {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.source_svg
}
func (t *Thumbs) set_filepath_source_img(fpath string) *Thumbs {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.source_img = typ_.FilePath(fpath)
	return t
}
func (t *Thumbs) get_filepath_source_img() typ_.FilePath {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.source_img
}
.....
*/

func (t *Thumbs) handle() {
	http.Handle(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		fpath, exists, err := t.ServeFile(r.URL)
		if err != nil {
			log.Println(t.l.Error(logT01, err))
			w.WriteHeader(http.StatusInternalServerError)
		} else if !exists {
			w.WriteHeader(http.StatusNotFound)
		} else {
			http.ServeFile(w, r, fpath)
		}
	}))
}

func (t *Thumbs) serve_file(url_ *url.URL) (fpath string, exists bool, err error) {

	if manifest, exists, err := t.server_file_manifest(url_); err != nil {
		return ``, false, t.l.Error(logTP, logT02, err)
	} else if exists {
		return manifest.String(), true, nil
	}

	if thumb, exists, err := t.server_file_thumb(url_, conv); err != nil {
		return ``, false, t.l.Error(logTP, logT03, err)
	} else if exists {
		return thumb.String(), true, nil
	}
	return ``, false, nil
}

func (t *Thumbs) server_file_thumb(url_ *url.URL, conv convert_.Converter) (fpath typ_.FilePath, exists bool, err error) {

	t.mu.RLock()
	thumb, exists := thumb_.URLExists(url_, t.thumbs)
	t.mu.RUnlock()

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
		return ``, false, t.l.Error(logTP, logT04, err)
	}

	return fpath, true, nil
}

func (t *Thumbs) server_file_manifest(url_ *url.URL) (manifest typ_.FilePath, exists bool, err error) {

	if !t.manifest.URLExists(url_) {
		return ``, false, nil
	}

	manifest, exists, err = t.manifest.GetFile(t.get_folder_work(), t.thumbs)
	if err != nil {
		return ``, false, t.l.Error(logTP, logT05, err)
	}
	return manifest, exists, nil
}

// func (T *Thumbs) manifest_url_exists( URLpath string ) ( fpath string, exists bool, err error )
// func (t *Thumbs) get_thumbs() map[types_.URLHref/*clear*/]*thumb_.Thumb

func (t *Thumbs) tags_html() string {

	var tags strings.Builder

	t.mu.RLock()
	for _, thumb := range t.thumbs {
		tags.WriteString(thumb.GetTAG())
	}
	tags.WriteString(t.manifest.GetTAG())
	t.mu.RUnlock()

	return tags.String()
}

func (t *Thumbs) getThumb(u typ_.URLHref) (*thumb_.Thumb, bool) {
	tb, ok := t.thumbs[u]
	return tb, ok
}

func default_list() *Thumbs {
	return &Thumbs{
		thumbs: func() map[typ_.URLHref]*thumb_.Thumb {
			m := map[typ_.URLHref]*thumb_.Thumb{}
			for _, thumb := range defaults_.Defaults() {
				m[thumb.GetHREF()] = thumb
			}
			return m
		}(),
		manifest: manifest_.Default(),
	}
}

func custom_list(thumbs ...*thumb_.Thumb) *Thumbs {
	return &Thumbs{
		thumbs: func() map[typ_.URLHref]*thumb_.Thumb {
			m := map[typ_.URLHref]*thumb_.Thumb{}
			for _, thumb := range thumbs {
				m[thumb.GetHREF()] = thumb
			}
			return m
		}(),
	}
}
