package thumbs

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 09 March 2023
 */
import (
	"io"
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
	logT07 = `T07: get new thumb`
)

var (
	TypePNG = thumb_.PNG
	TypeICO = thumb_.ICO
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

func (t *Thumbs) NewThumb(key string, typThumb thumb_.Typ) (*thumb_.Thumb, error) {
	tb, err := thumb_.NewThumb(key, typThumb, t.l, t.storage, t.conv)
	if err != nil {
		return nil, t.l.Error(logTP, logT07, err)
	}
	return tb, nil
}

func (t *Thumbs) AppendThumb(tb *thumb_.Thumb) *Thumbs {
	return t.append(tb)
}



// получение конкретной превьюхи для отправки пользователю
func (t *Thumbs) File(urlPath string) (content io.ReadSeekCloser, modtime time.Time, name string, exists bool, err error) {
	return t.thumbFile(urlPath)
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
