package thumb

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	"html"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	create_ "github.com/jhekau/favicon/thumb/create"
	types_ "github.com/jhekau/favicon/types"
)

var (
	// ~~ interface ~~


	URLExists = url_Exists
)

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


///
///
type Thumb struct {
	sync.RWMutex
	size_px uint16
	size_attr_value attr_size
	comment string // <!-- comment -->
	url_href types_.URLHref // domain{/name_url}, first -> `/`
	url_href_clear types_.URLHref 
	filename string // [folder/file] [file] [.file]
	tag_rel string
	manifest bool
	mimetype types_.FileType
	typ types_.FileType
	cache *cache
}

//
// ...
func (t *Thumb) file_create(save_img, source_img, source_svg types_.FilePath) (complite bool, err error) {

	t.Lock()
	defer t.Unlock()

	os.Remove(save_img.String())
	complite, err = create_.File(source_img, source_svg, save_img, t.typ, int(t.size_px))
	if err != nil {
		// return error
	}
	return complite, nil
}

// ...
func (t *Thumb) get_filepath(folder_work types_.Folder, original_name types_.FileName) types_.FilePath {

	var fpath types_.FilePath
	t.RLock()
	fpath = t.cache.get_filepath()
	t.RUnlock()

	if fpath != `` {
		return fpath
	}

	//
	if t.filename == `` {
		t.filename = strings.Join(
			[]string{
				original_name.String(),
				strconv.Itoa(int(t.size_px)),
			},
			`_`,
		) + `.` + t.mimetype.String()
	}
	
	fpath = types_.FilePath(
		filepath.Join(folder_work.String(), t.filename),
	)

	//
	t.cache.set_filepath(fpath)
	return fpath
}

// ...
func (t *Thumb) GetFile(folder_work types_.Folder, source_img, source_svg types_.FilePath) (types_.FilePath, error) {

	original_filename := types_.FileName(
		filepath.Base(source_img.String()),
	)

	save_img := t.get_filepath(folder_work, original_filename)
	var check_exists types_.FileExists

	t.RLock()
	check_exists = t.cache.get_file_exists_state()
	t.RUnlock()

	if check_exists == types_.FileExistsOK {
		return save_img, nil
	}
	
	if check_exists == types_.FileExistsNOT {
		// return error
	}

	t.Lock()
	defer t.Unlock()


	if complite, err := t.file_create(save_img, source_img, source_svg); err != nil {
		// return error
	} else if !complite {
		t.cache.set_file_exists_state(types_.FileExistsNOT)
		// return error
	}

	if f, err := os.Stat(save_img.String()); err != nil {
		t.cache.set_file_exists_state(types_.FileExistsNOT)
		if os.IsNotExist(err) {
			// return ``, error - file not exists
		}
		// return error
	} else if f.IsDir() {
		t.cache.set_file_exists_state(types_.FileExistsNOT)
		// return error - file is folder
	}

	t.cache.set_file_exists_state(types_.FileExistsOK)
	return save_img, nil
}

// ...
func (t *Thumb) SetSize(px uint16) *Thumb {

	t.Lock()
	defer t.Unlock()

	t.cache.clean()
	t.size_px = px
	return t
}

func (t *Thumb) GetSize() uint16 {

	t.RLock()
	defer t.RUnlock()

	return t.size_px
}

// ...
// func (t *Thumb) SetNameURL( nameURL string ) *Thumb {
// 	t.cache.clean()
// 	t.url_path = nameURL
// 	return t
// }

// ...
func (t *Thumb) SetNameFile( nameFile string ) *Thumb {

	t.Lock()
	defer t.Unlock()

	t.cache.clean()
	t.filename = nameFile
	return t
}

// ...
func (t *Thumb) SetTagRel( tagRel string ) *Thumb {

	t.Lock()
	defer t.Unlock()

	t.cache.clean()
	t.tag_rel = tagRel
	return t
}

// ...
func (t *Thumb) SetManifestUsed() *Thumb {

	t.Lock()
	defer t.Unlock()

	t.cache.clean()
	t.manifest = true
	return t
}

// <!-- comment -->
func (t *Thumb) SetHTMLComment(comment string) {

	t.Lock()
	defer t.Unlock()

	t.comment = comment
}

// ...
func (t *Thumb) SetType(typ types_.FileType) *Thumb {

	t.Lock()
	defer t.Unlock()

	t.cache.clean()
	t.typ = typ
	return t
}

func (t *Thumb) GetType() types_.FileType {

	t.RLock()
	defer t.RUnlock()

	return t.typ
}

// ...
func (t *Thumb) SetHREF(src string) *Thumb {

	t.Lock()
	defer t.Unlock()

	t.cache.clean()
	t.url_href = types_.URLHref(src)
	{
		u, err := url.Parse(`http://domain.com`)
		if err != nil {
			// error
		} else {
			t.url_href_clear = types_.URLHref(u.JoinPath(src).Path)
		}
	}
	return t
}

func (t *Thumb) GetHREF() types_.URLHref {

	t.RLock()
	defer t.RUnlock()

	return t.url_href
}
func (t *Thumb) GetHREFClear() types_.URLHref {

	t.RLock()
	defer t.RUnlock()

	return t.url_href_clear
}



// ...
func (t *Thumb) StatusManifest() bool { // ( string, bool /*true - used*/ )

	t.RLock()
	defer t.RUnlock()

	return t.manifest
}

// ...
func (t *Thumb) GetTAG() string {

	t.RLock()
	if str := t.cache.get_tag(); str != `` {
		t.Unlock()
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

	t.RUnlock()

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
		t.Lock()
		t.cache.set_tag(str)
		t.Unlock()
	}

	return str
}

// ...
func (t *Thumb) SetTypeImage( typ types_.FileType ) *Thumb {
	
	t.Lock()
	defer t.Unlock()

	t.cache.clean()
	t.typ = typ
	return t
}

// ...
func (t *Thumb) SetSizeAttrEmpty() *Thumb {

	t.Lock()
	defer t.Unlock()

	t.cache.clean()
	t.size_attr_value = attr_size{
		state: attr_size_state_empty,
	}
	return t
}

func (t *Thumb) SetSizeAttrDefault() *Thumb {

	t.Lock()
	defer t.Unlock()

	t.cache.clean()
	t.size_attr_value = attr_size{
		state: attr_size_state_custom,
	}
	return t
}

func (t *Thumb) SetSizeAttrCustom(val string) *Thumb {

	t.Lock()
	defer t.Unlock()

	t.cache.clean()
	t.size_attr_value = attr_size{
		state: attr_size_state_custom,
		value: val,
	}
	return t
}














type cache struct {
	sync.RWMutex
	filepath types_.FilePath
	tag string
	file_exists_state types_.FileExists
}

//
func (c *cache) get_filepath() types_.FilePath {

	c.RLock()
	defer c.RUnlock()

	if c != nil {
		return c.filepath
	}
	return ``
}
func (c *cache) set_filepath( filepath types_.FilePath ) {
	c.Lock()
	if c == nil {
		c = &cache{}
	}
	c.filepath = filepath
	c.Unlock()
}

//
func (c *cache) get_tag() string {

	c.RLock()
	defer c.RUnlock()

	if c != nil {
		return c.tag
	}
	return ``
}
func (c *cache) set_tag( tag string ) {
	c.Lock()
	if c == nil {
		c = &cache{}
	}
	c.tag = tag
	c.Unlock()
}

//
func (c *cache) get_file_exists_state() types_.FileExists {

	c.RLock()
	defer c.RUnlock()

	if c != nil {
		return c.file_exists_state
	}
	return types_.FileExistsNotCheck
}

func (c *cache) set_file_exists_state( exists types_.FileExists ) {
	c.Lock()
	if c == nil {
		c = &cache{}
	}
	c.file_exists_state = exists
	c.Unlock()
}

//
func (c *cache) clean() {
	c.Lock()
	c = nil
	c.Unlock()
}






// url_Exists : проверка наличия превьюхи в запросе 
// http.Request.URL.Path -> URLpath
func url_Exists( URLpath string, thumbs map[string /*url_href_clear*/]*Thumb ) ( *Thumb, bool /*exists*/ ) {
	t, ok := thumbs[URLpath]
	return t, ok
}

