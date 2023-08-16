package thumb

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	"html"
	"io"
	"time"

	"strconv"
	"sync"

	typ_ "github.com/jhekau/favicon/internal/core/types"
	types_ "github.com/jhekau/favicon/pkg/core/types"
	converter_ "github.com/jhekau/favicon/pkg/core/models/converter"
	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	err_ "github.com/jhekau/favicon/internal/core/err"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"
)

const (
	logTP  = `/thumb/thumb.go`
	logT01 = `T01: create file`
	logT02 = `T02: create new object thumb into storage`
	logT03 = `T03: original image undefined`
	// logT04 = `T04: `
	// logT05 = `T05: `
	// logT06 = `T06: `
	// logT07 = `T07: `
	// logT08 = `T08: `
	// logT09 = `T09: create new storage object`

	logT10 = `T10: thumb check is exists`
	logT11 = `T11: create new thumb image`
	// logT12 = `T12: `
	// logT13 = `T13: `
	// logT14 = `T14: `
	// logT15 = `T15: `
	// logT16 = `T16: `
	// logT17 = `T17: `
	// logT18 = `T18: `
	// logT19 = `T19: `
)

var (
	URLPath_Get = urlPath_Get
)

type Typ types_.FileType
var (
	ICO Typ = Typ(types_.ICO())
	PNG Typ = Typ(types_.PNG())
	SVG Typ = Typ(types_.SVG())
)


type cache interface{
	Delete(key any)
	Load(key any) (value any, ok bool)
	Range(f func(key any, value any) bool)
	Store(key any, value any)
}


type attrSizeState int
const (
	attrSizeEmpty = iota-1
	attrSizeDefault
	attrSizeCustom
)
type attrSize struct {
	state attrSizeState
	val string
}

// Оригинальное изображение, с которого нарезается превьюха
type original struct {
	typSVG bool
	obj storage_.StorageOBJ
}

// внимание! ключь, если используется файловая система в качестве
// хранилища по умолчанию, используется как filepath для хранения превьюхи
func NewThumb(key string, typThumb Typ, l logger_.Logger, s storage_.Storage, c converter_.Converter) (*Thumb, error) {
	t, err := s.NewObject(key)
	if err != nil {
		return nil, err_.Err(l, logTP, logT02, err)
	}
	return &Thumb{
		l:l,
		storage:s,
		conv:c, 
		cache: &sync.Map{},
		thumb: t,
		mimetype: types_.FileType(typThumb),
	}, nil
}

type Thumb struct {

	mu sync.RWMutex

	l logger_.Logger
	storage storage_.Storage
	conv converter_.Converter
	cache cache // from test

	original *original
	thumb storage_.StorageOBJ

	attrSize attrSize
	attrRel string

	sizePX int
	comment string // <!-- comment -->
	urlPath typ_.URLPath // domain{/name_url}, first -> `/`
	manifest bool
	mimetype types_.FileType
	
}

func (t *Thumb) SetSize(px int) *Thumb {
	return t.setSize(px)
}

func (t *Thumb) GetSize() int {
	return t.getSize()
}

func (t *Thumb) SetAttrRel( tagRel string ) *Thumb {
	return t.setAttrRel(tagRel)
}

// Добавлять ли превью в список манифеста
func (t *Thumb) SetManifestUsed() *Thumb {
	return t.setManifestUsed()
}

func (t *Thumb) SetHTMLComment(comment string) *Thumb {
	return t.setHTMLComment(comment)
}

func (t *Thumb) GetType() types_.FileType {
	return t.getType()
}

func (t *Thumb) SetURLPath(src string) *Thumb {
	return t.setURLPath(src)
}

func (t *Thumb) GetURLPath() typ_.URLPath {
	return t.getURLPath()
}

func (t *Thumb) StatusManifest() bool { // ( string, bool /*true - used*/ )
	return t.statusManifest()
}

func (t *Thumb) GetTAG() string {
	return t.getTag()
}

// аттрибут size не будет добавлен в тег
func (t *Thumb) SetAttrSize_Empty() *Thumb {
	return t.setAttrSize_Empty()
}

// аттрибут size будет добавлен только в том случае, если указан размер превью
func (t *Thumb) SetAttrSize_Default() *Thumb {
	return t.setAttrSize_Default()
}

// аттрибут size будет содержать кастомное значение val
func (t *Thumb) SetAttrSize_Custom(val string) *Thumb {
	return t.setAttrSize_Custom(val)
}

func (t *Thumb) GetOriginalKey() string{
	return t.getOriginalKEY()
}

func (t *Thumb) Read() (io.ReadSeekCloser, error) {
	return t.read()
}

func (t *Thumb) ModTime() time.Time {
	return t.modtime()
}



// добавление оригинального изображения для нарезки превьюхи
func (t *Thumb) SetOriginal( obj storage_.StorageOBJ ) *Thumb {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.original = &original{
		obj: obj,
	}
	return t
}

// добавление оригинального изображения SVG для нарезки превьюхи
func (t *Thumb) SetOriginalSVG( obj storage_.StorageOBJ ) *Thumb {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.original = &original{
		typSVG: true,
		obj: obj,
	}
	return t
}



func (t *Thumb) getOriginal() *original {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.original
}


func (t *Thumb) modtime() time.Time {
	t.mu.RLock()
	tt := t.thumb
	t.mu.RUnlock()

	if tt == nil {
		return time.Time{}
	}
	return tt.ModTime()
}


func (t *Thumb) thumb_create() error {

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.original == nil {
		return err_.Err(t.l, logTP, logT03)
	}

	err := t.conv.Do(t.original.obj, t.thumb, t.original.typSVG, t.mimetype, int(t.sizePX))
	if err != nil {
		return err_.Err(t.l, logTP, logT01, err)
	}
	return nil
}

func (t *Thumb) getOriginalKEY() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.original == nil {
		return ``
	}
	return string(t.original.obj.Key())
}

func (t *Thumb) read() (io.ReadSeekCloser, error) {
	
	t.mu.RLock()
	thumb := t.thumb
	t.mu.RUnlock()
	
	exist, err := thumb.IsExists()
	if err != nil {
		return nil, err_.Err(t.l, logTP, logT10, err)
	}
	if exist {
		return thumb.Reader()
	}

	err = t.thumb_create()
	if err != nil {
		return nil, err_.Err(t.l, logTP, logT11, err)
	}
	return t.thumb.Reader()
}

// ...
func (t *Thumb) setSize(px int) *Thumb {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.cacheClean()
	t.sizePX = px
	return t
}

func (t *Thumb) getSize() int {

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.sizePX
}

// ...
func (t *Thumb) setAttrRel( tagRel string ) *Thumb {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.cacheClean()
	t.attrRel = tagRel
	return t
}

// ...
func (t *Thumb) setManifestUsed() *Thumb {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.cacheClean()
	t.manifest = true
	return t
}

// <!-- comment -->
func (t *Thumb) setHTMLComment(comment string) *Thumb {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.comment = comment
	return t
}

func (t *Thumb) getType() types_.FileType {

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.mimetype
}

// ...
func (t *Thumb) setURLPath(src string) *Thumb {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.cacheClean()
	t.urlPath = typ_.URLPath(src)

	return t
}

func (t *Thumb) getURLPath() typ_.URLPath {

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.urlPath
}


// ...
func (t *Thumb) statusManifest() bool { // ( string, bool /*true - used*/ )

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.manifest
}

// ...
func (t *Thumb) getTag() string {

	t.mu.RLock()
	defer t.mu.RUnlock()

	if str := t.tagCacheGet(); str != `` {
		return str
	}

	// <link rel="apple-touch-icon" sizes="180x180" href="touch-icon-iphone-retina.png" type="image/png">

	attr := map[string]string{}

	// size
	switch t.attrSize.state {
	case attrSizeEmpty:
	case attrSizeDefault:
		sz := strconv.Itoa(int(t.sizePX))
		if t.sizePX > 0 {
			attr[`sizes`] = sz+`x`+sz
		}
	case attrSizeCustom:
		attr[`sizes`] = html.EscapeString(t.attrSize.val)
	}

	// href
	if s := t.urlPath.String(); s != `` {
		attr[`href`] = html.EscapeString(s)
	}

	// rel
	if t.attrRel != `` {
		attr[`rel`] = html.EscapeString(t.attrRel)
	}
	
	// type
	if t.mimetype != `` {
		attr[`type`] = html.EscapeString(t.mimetype.String())
	}

	// if comment <tag /> <!-- comment -->
	comment := t.comment

	str := ``
	if len(attr) > 0 {

		str += `<link`
		for name, val := range attr {
			str += ` `+name+`="`+val+`" `
		}
		str += `>`

		if comment != `` {
			str += `<!-- `+html.EscapeString(comment)+` -->`
		}
		t.tagCacheSet(str)
	}

	return str
}

func (t *Thumb) tagCacheGet() string {
	c, ok := t.cache.Load(`tag`)
	if ok {
		return c.(string)
	}
	return ``
}
func (t *Thumb) tagCacheSet( s string ) {
	t.cache.Store(`tag`, s)
}

// ...
func (t *Thumb) setAttrSize_Empty() *Thumb {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.cacheClean()
	t.attrSize = attrSize{
		state: attrSizeEmpty,
	}
	return t
}

func (t *Thumb) setAttrSize_Default() *Thumb {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.cacheClean()
	t.attrSize = attrSize{
		state: attrSizeDefault,
	}
	return t
}

func (t *Thumb) setAttrSize_Custom(val string) *Thumb {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.cacheClean()
	t.attrSize = attrSize{
		state: attrSizeCustom,
		val: val,
	}
	return t
}



func (t *Thumb) cacheClean() {
	t.cache.Range(func(k, _ any) bool{
		t.cache.Delete(k)
		return true
	})
}









func urlPath_Get( urlPath string, thumbs map[typ_.URLPath]*Thumb ) ( *Thumb, bool /*exists*/ ) {
	t, ok := thumbs[typ_.URLPath(urlPath)]
	return t, ok
}

