package thumb

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	"html"
	"io"
	"log"
	"net/url"
	"strconv"
	"sync"

	logger_ "github.com/jhekau/favicon/internal/core/logger"
	types_ "github.com/jhekau/favicon/internal/core/types"
	files_ "github.com/jhekau/favicon/internal/storage/files"
	domain_ "github.com/jhekau/favicon/pkg/domain"
)

const (
	logTP  = `/thumb/thumb.go`
	logT01 = `T01: create file`
	// logT02 = `T02: thumb file not exists`
	// logT03 = `T03: create thumb file`
	// logT04 = `T04: not complite - file create`
	// logT05 = `T05: thumb not exists`
	// logT06 = `T06: os stat save thumb`
	// logT07 = `T07: save thumb is a folder`
	logT08 = `T08: url parse standart template domain.com`
	logT09 = `T09: create new storage object`

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
	URLExists = url_Exists
)


type StorageOBJ interface{
	Reader() (io.ReadCloser , error)
	Writer() (io.WriteCloser, error)
	Key() domain_.StorageKey
	IsExists() ( bool, error )
}

type Storage interface {
	NewObject() (StorageOBJ, error)
}

type Converter interface{
	Do(source, save StorageOBJ, originalSVG bool, typThumb types_.FileType, size_px int) error
}

// type Converter interface{
// 	Do(source, save StorageOBJ, originalSVG bool, typ types_.FileType, size_px int) error
// }


///
///
type attr_size_state int

const (
	attr_size_state_empty attr_size_state = -1
	attr_size_state_default attr_size_state = 0
	attr_size_state_custom attr_size_state = 1
)

type attr_size struct {
	state attr_size_state
	value string
}

// Оригинальное изображение, с которого нарезается превьюха
type original struct {
	typSVG bool
	// filepath types_.FilePath
	obj StorageOBJ
}

///
///
type Thumb struct {

	s sync.RWMutex
	l *logger_.Logger

	original *original
	thumb StorageOBJ

	storage Storage
	conv Converter

	size_px uint16
	size_attr_value attr_size
	comment string // <!-- comment -->
	url_href types_.URLHref // domain{/name_url}, first -> `/`
	url_href_clear types_.URLHref 
	// filename string 
	tag_rel string
	manifest bool
	mimetype types_.FileType
	typ types_.FileType
	cache cache

}

func (t *Thumb) SetSize(px uint16) *Thumb {
	return t.set_size(px)
}

func (t *Thumb) GetSize() uint16 {
	return t.get_size()
}

// func (t *Thumb) SetNameFile( nameFile string ) *Thumb {
// 	return t.set_name_file(nameFile)
// }

func (t *Thumb) SetTagRel( tagRel string ) *Thumb {
	return t.set_tag_rel(tagRel)
}

func (t *Thumb) SetManifestUsed() *Thumb {
	return t.set_manifest_used()
}

func (t *Thumb) SetHTMLComment(comment string) *Thumb {
	return t.set_html_comment(comment)
}

func (t *Thumb) SetType(typ types_.FileType) *Thumb {
	return t.set_type(typ)
}

func (t *Thumb) GetType() types_.FileType {
	return t.get_type()
}

func (t *Thumb) SetHREF(src string) *Thumb {
	return t.set_href(src)
}

func (t *Thumb) GetHREF() types_.URLHref {
	return t.get_href()
}

func (t *Thumb) GetHREFClear() types_.URLHref {
	return t.get_href_clear()
}

func (t *Thumb) StatusManifest() bool { // ( string, bool /*true - used*/ )
	return t.status_manifest()
}

func (t *Thumb) GetTAG() string {
	return t.get_tag()
}

// func (t *Thumb) SetTypeImage( typ types_.FileType ) *Thumb {
// 	return t.set_type_image(typ)
// }

func (t *Thumb) SetSizeAttrEmpty() *Thumb {
	return t.set_size_attr_empty()
}

func (t *Thumb) SetSizeAttrDefault() *Thumb {
	return t.set_size_attr_default()
}

func (t *Thumb) SetSizeAttrCustom(val string) *Thumb {
	return t.set_size_attr_custom(val)
}

func (t *Thumb) GetOriginalKey() string{
	return t.get_original_key()
}

func (t *Thumb) Read() (io.ReadCloser, error) {
	return t.read()
}




// source image
func (t *Thumb) OriginalFileSet( filepath string ) {
	file := (&files_.Files{L: t.l}).NewObject(types_.FilePath(filepath))
	t.original = &original{
		obj: file,
	}
}
func (t *Thumb) OriginalFileSetSVG( filepath string ) {
	file := (&files_.Files{L: t.l}).NewObject(types_.FilePath(filepath))
	t.original = &original{
		typSVG: true,
		obj: file,
	}
}
func (t *Thumb) OriginalCustomSet( obj StorageOBJ ) {
	t.original = &original{
		obj: obj,
	}
}
func (t *Thumb) OriginalCustomSetSVG( obj StorageOBJ ) {
	t.original = &original{
		typSVG: true,
		obj: obj,
	}
}



func (t *Thumb) original_get( filepath string ) *original {
	return t.original
}





func (t *Thumb) thumb_create() error {

	t.s.Lock()
	defer t.s.Unlock()

	err := t.conv.Do(t.original.obj, t.thumb, t.original.typSVG, t.typ, int(t.size_px))
	if err != nil {
		return t.l.Typ.Error(logTP, logT01, err)
	}
	return nil
}


func (t *Thumb) get_original_key() string {
	if t.original == nil {
		return ``
	}
	return string(t.original.obj.Key())
}

func (t *Thumb) read() (io.ReadCloser, error) {
	
	if t.thumb == nil {
		tb, err := t.storage.NewObject()
		if err != nil {
			return nil, t.l.Typ.Error(logTP, logT09, err)
		}
		t.thumb = tb
	}
	exist, err := t.thumb.IsExists()
	if err != nil {
		return nil, t.l.Typ.Error(logTP, logT10, err)
	}
	if exist {
		return t.thumb.Reader()
	}

	err = t.thumb_create()
	if err != nil {
		return nil, t.l.Typ.Error(logTP, logT11, err)
	}
	return t.thumb.Reader()
}

// ...
func (t *Thumb) set_size(px uint16) *Thumb {

	t.s.Lock()
	defer t.s.Unlock()

	t.cache.clean()
	t.size_px = px
	return t
}

func (t *Thumb) get_size() uint16 {

	t.s.RLock()
	defer t.s.RUnlock()

	return t.size_px
}

// ...
// func (t *Thumb) SetNameURL( nameURL string ) *Thumb {
// 	t.cache.clean()
// 	t.url_path = nameURL
// 	return t
// }

// ...
func (t *Thumb) set_tag_rel( tagRel string ) *Thumb {

	t.s.Lock()
	defer t.s.Unlock()

	t.cache.clean()
	t.tag_rel = tagRel
	return t
}

// ...
func (t *Thumb) set_manifest_used() *Thumb {

	t.s.Lock()
	defer t.s.Unlock()

	t.cache.clean()
	t.manifest = true
	return t
}

// <!-- comment -->
func (t *Thumb) set_html_comment(comment string) *Thumb {

	t.s.Lock()
	defer t.s.Unlock()

	t.comment = comment
	return t
}

// ...
func (t *Thumb) set_type(typ types_.FileType) *Thumb {

	t.s.Lock()
	defer t.s.Unlock()

	t.cache.clean()
	t.typ = typ
	return t
}

func (t *Thumb) get_type() types_.FileType {

	t.s.RLock()
	defer t.s.RUnlock()

	return t.typ
}

// ...
func (t *Thumb) set_href(src string) *Thumb {

	t.s.Lock()
	defer t.s.Unlock()

	t.cache.clean()
	t.url_href = types_.URLHref(src)
	{
		u, err := url.Parse(`http://domain.com`)
		if err != nil {
			log.Println( t.l.Typ.Error(logTP, logT08, err) )
		} else {
			t.url_href_clear = types_.URLHref(u.JoinPath(src).Path)
		}
	}
	return t
}

func (t *Thumb) get_href() types_.URLHref {

	t.s.RLock()
	defer t.s.RUnlock()

	return t.url_href
}
func (t *Thumb) get_href_clear() types_.URLHref {

	t.s.RLock()
	defer t.s.RUnlock()

	return t.url_href_clear
}



// ...
func (t *Thumb) status_manifest() bool { // ( string, bool /*true - used*/ )

	t.s.RLock()
	defer t.s.RUnlock()

	return t.manifest
}

// ...
func (t *Thumb) get_tag() string {

	t.s.RLock()
	if str := t.cache.get_tag(); str != `` {
		t.s.Unlock()
		return str
	}

	// <link rel="apple-touch-icon" sizes="180x180" href="touch-icon-iphone-retina.png" type="image/png">

	attr := map[string]string{}

	// size
	switch t.size_attr_value.state {
	case attr_size_state_empty:
	case attr_size_state_default:
		sz := strconv.Itoa(int(t.size_px))
		attr[`sizes`] = sz+`x`+sz
	case attr_size_state_custom:
		attr[`sizes`] = html.EscapeString(t.size_attr_value.value)
	}

	// href
	if s := t.url_href.String(); s != `` {
		attr[`href`] = html.EscapeString(s)
	}

	// rel
	if t.tag_rel != `` {
		attr[`rel`] = html.EscapeString(t.tag_rel)
	}
	
	// type
	if t.mimetype != `` {
		attr[`type`] = html.EscapeString(t.mimetype.String())
	}

	// if comment <tag /> <!-- comment -->
	comment := t.comment

	t.s.RUnlock()

	str := ``
	if len(attr) > 0 {

		str += `<`
		for name, val := range attr {
			str += ` `+name+`="`+val+`" `
		}
		str += `>`

		if comment != `` {
			str += `<!-- `+html.EscapeString(comment)+` -->`
		}
		t.s.Lock()
		t.cache.set_tag(str)
		t.s.Unlock()
	}

	return str
}

// ...
func (t *Thumb) set_size_attr_empty() *Thumb {

	t.s.Lock()
	defer t.s.Unlock()

	t.cache.clean()
	t.size_attr_value = attr_size{
		state: attr_size_state_empty,
	}
	return t
}

func (t *Thumb) set_size_attr_default() *Thumb {

	t.s.Lock()
	defer t.s.Unlock()

	t.cache.clean()
	t.size_attr_value = attr_size{
		state: attr_size_state_custom,
	}
	return t
}

func (t *Thumb) set_size_attr_custom(val string) *Thumb {

	t.s.Lock()
	defer t.s.Unlock()

	t.cache.clean()
	t.size_attr_value = attr_size{
		state: attr_size_state_custom,
		value: val,
	}
	return t
}














type cache struct {
	sync.RWMutex
	tag string
}

//
func (c *cache) get_tag() string {

	c.RLock()
	defer c.RUnlock()

	return c.tag
}
func (c *cache) set_tag( tag string ) {
	c.Lock()
	c.tag = tag
	c.Unlock()
}

//
func (c *cache) clean() {
	c.Lock()
	c = &cache{}
	c.Unlock()
}






// url_Exists : проверка наличия превьюхи в запросе 
// http.Request.URL.Path -> URLpath
func url_Exists( url_ *url.URL, thumbs map[types_.URLHref /*url_href_clear*/]*Thumb ) ( *Thumb, bool /*exists*/ ) {
	t, ok := thumbs[types_.URLHref(url_.Path)]
	return t, ok
}

