package thumbs

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 09 March 2023
 */
import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

// создание набора превьюх по умолчанию для оригинала
func NewThumbs_DefaultsIcons() *Thumbs {
	t := NewThumbs()
	return t
}

type original struct {
	image storage_.StorageOBJ
	svg bool
}

type Thumbs struct {
	mu          sync.RWMutex
	l           logger_.Logger
	storage     storage_.Storage
	conv 		converter_.Converter
	original 	*original
	folder_work typ_.Folder
	thumbs      map[typ_.URLPath]*thumb_.Thumb
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

func (t *Thumbs) Handle() {
	t.handle()
}




func (t *Thumbs) TagsHTML() string {
	return t.tags_html()
}

func (t *Thumbs) append(thumb *thumb_.Thumb) *Thumbs {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.thumbs[thumb.URLPathGet()] = thumb
	return t
}
/*
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
*/

func (t *Thumbs) handle() {
	http.Handle(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		content, modtime, name, exists, err := t.thumbFile(r.URL.Path)
		if err != nil {
			log.Println(t.l.Error(logT01, err))
			w.WriteHeader(http.StatusInternalServerError)
		} else if !exists {
			w.WriteHeader(http.StatusNotFound)
		} else {
			http.ServeContent(w, r, name, modtime, content)
		}
	}))
}

// func (t *Thumbs) contentFile(url_ *url.URL) (content io.ReadSeekCloser, modtime time.Time, name string, exists bool, err error) {

// 	if c,tm,n,e,er := t.thumbFile(url_); err != nil {
// 		err = t.l.Error(logTP, logT03, er)
// 		return
// 	} else if exists {
// 		return c,tm,n,e,err
// 	}
// 	return ``, false, nil
// }

func (t *Thumbs) thumbFile(urlPath string) (content io.ReadSeekCloser, modtime time.Time, name string, exists bool, err error) {

	t.mu.RLock()
	thumb, exists := thumb_.URLPath_Get(urlPath, t.thumbs)
	t.mu.RUnlock()

	if !exists {
		return
	}
	modtime = thumb.ModTime()
	_, name = filepath.Split(urlPath)

	content, err = thumb.Read()
	if err != nil {
		err = t.l.Error(logTP, logT04, err)
	}
	return
}

// func (t *Thumbs) server_file_manifest(url_ *url.URL) (manifest typ_.FilePath, exists bool, err error) {

// 	if !t.manifest.URLExists(url_) {
// 		return ``, false, nil
// 	}

// 	manifest, exists, err = t.manifest.GetFile(t.get_folder_work(), t.thumbs)
// 	if err != nil {
// 		return ``, false, t.l.Error(logTP, logT05, err)
// 	}
// 	return manifest, exists, nil
// }

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

func (t *Thumbs) getThumb(u typ_.URLPath) (*thumb_.Thumb, bool) {
	tb, ok := t.thumbs[u]
	return tb, ok
}

func default_list() *Thumbs {
	return &Thumbs{
		thumbs: func() map[typ_.URLPath]*thumb_.Thumb {
			m := map[typ_.URLPath]*thumb_.Thumb{}
			for _, thumb := range defaults_.Defaults() {
				m[thumb.URLPathGet()] = thumb
			}
			return m
		}(),
		manifest: manifest_.Default(),
	}
}

func custom_list(thumbs ...*thumb_.Thumb) *Thumbs {
	return &Thumbs{
		thumbs: func() map[typ_.URLPath]*thumb_.Thumb {
			m := map[typ_.URLPath]*thumb_.Thumb{}
			for _, thumb := range thumbs {
				m[thumb.URLPathGet()] = thumb
			}
			return m
		}(),
	}
}
