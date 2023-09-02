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

	logs_ "github.com/jhekau/favicon/internal/pkg/logs"
	typ_ "github.com/jhekau/favicon/internal/pkg/types"
	convert_ "github.com/jhekau/favicon/internal/pkg/img/convert"
	checks_ "github.com/jhekau/favicon/internal/pkg/img/convert/checks"
	converters_ "github.com/jhekau/favicon/internal/pkg/img/convert/converters"
	defaults_ "github.com/jhekau/favicon/internal/service/thumb/defaults"
	converter_exec_anthonynsimon_ "github.com/jhekau/favicon/internal/pkg/img/convert.exec/anthonynsimon"
	resolution_ "github.com/jhekau/favicon/internal/pkg/img/resolution"
	manifest_ "github.com/jhekau/favicon/internal/service/manifest"
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
	files_ "github.com/jhekau/favicon/internal/storage/files"

	converter_ "github.com/jhekau/favicon/interfaces/converter"
	logger_ "github.com/jhekau/favicon/interfaces/logger"
	storage_ "github.com/jhekau/favicon/interfaces/storage"
	err_ "github.com/jhekau/favicon/internal/pkg/err"
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
	logT08 = `T08: get list defaults`
	logT09 = `T09: original empty`
)

var (
	TypePNG = thumb_.PNG
	TypeICO = thumb_.ICO
	TypeSVG = thumb_.SVG
)

// создание пустого набора превьюх для одного оригинального изображения
func NewThumbs() *Thumbs {
	logger := &logs_.Logger{}
	return &Thumbs{
		l: logger,
		storage: files_.NewStorage(``, logger).SetDirDefault(),
		conv: &convert_.Converter{
			L: logger,
			Converters: []converter_.ConverterTyp{
				&converters_.ConverterICO{
					ConverterExec: &converter_exec_anthonynsimon_.Exec{L: logger},
				},
				&converters_.ConverterPNG{
					ConverterExec: &converter_exec_anthonynsimon_.Exec{L: logger},
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
		thumbs: make(map[typ_.URLPath]*thumb_.Thumb),
	}
}

// создание набора превьюх по умолчанию для оригинала
func NewThumbs_DefaultsIcons() (*Thumbs, error) {
	
	t := NewThumbs()

	icons, err := defaults_.Defaults(t.l, t.storage, t.conv)
	if err != nil {
		return nil, err_.Err(t.l, logTP, logT08, err)
	}

	for _, tb := range icons {
		t.AppendThumb(tb)
	}
	return t, nil
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
		return nil, err_.Err(t.l, logTP, logT07, err)
	}
	return tb, nil
}

func (t *Thumbs) AppendThumb(tb *thumb_.Thumb) *Thumbs {
	return t.append(tb)
}

func (t *Thumbs) SetOriginal(obj storage_.StorageOBJ) {
	t.original = &original{obj, false}
}

func (t *Thumbs) SetOriginalSVG(obj storage_.StorageOBJ) {
	t.original = &original{obj, true}
}

// получение конкретной превьюхи для отправки пользователю
func (t *Thumbs) ServeFile(urlPath string) (content io.ReadSeekCloser, modtime time.Time, name string, exists bool, err error) {
	return t.thumbFile(urlPath)
}




func (t *Thumbs) TagsHTML() string {
	return t.tags_html()
}

func (t *Thumbs) append(thumb *thumb_.Thumb) *Thumbs {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.thumbs[thumb.GetURLPath()] = thumb
	return t
}

func (t *Thumbs) thumbFile(urlPath string) (content io.ReadSeekCloser, modtime time.Time, name string, exists bool, err error) {

	t.mu.RLock()
	thumb, exists := thumb_.URLPath_Get(urlPath, t.thumbs)
	t.mu.RUnlock()

	if !exists {
		return
	}

	if thumb.GetOriginalKey() == `` {
		if t.original == nil {
			err = err_.Err(t.l, logTP, logT09)
			return
		}
		switch t.original.svg{
		case false: thumb.SetOriginal(t.original.image)
		case true: thumb.SetOriginalSVG(t.original.image)
		}
	}

	content, err = thumb.Read()
	if err != nil {
		err = err_.Err(t.l, logTP, logT04, err)
	}
	modtime = thumb.ModTime()
	name = filepath.Base(urlPath)

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
