package thumb

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
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
)

type Thumb struct {
	sync.RWMutex
	size_px int
	comment string // <!-- comment -->
	url_path string // domain{/name_url}, first -> `/`
	filename string // [folder/file] [file] [.file]
	tag_rel string
	manifest bool
	mimetype types_.FileType
	typ types_.FileType
	cache *cache
}

// ...
func (t *Thumb) file_create(save_img, source_img, source_svg types_.FilePath) (complite bool, err error) {
	os.Remove(save_img.String())
	complite, err = create_.File(source_img, source_svg, save_img, t.typ, t.size_px)
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
				strconv.Itoa(t.size_px),
				t.mimetype.String(),
			},
			`_`,
		)
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
func (t *Thumb) SetSize(px int) *Thumb {
	t.cache.clean()
}

func (t *Thumb) GetSize() int {
	return t.size_px
}

// ...
func (t *Thumb) SetNameURL( nameURL string ) *Thumb {
	t.cache.clean()
}

// ...
func (t *Thumb) SetNameFile( nameFile string ) *Thumb {
	t.cache.clean()
}

// ...
func (t *Thumb) SetTagRel( tagRel string ) *Thumb {
	t.cache.clean()
}

// ...
func (t *Thumb) SetManifestUsed() *Thumb {
	t.cache.clean()
}

// ...
func (t *Thumb) SetType(typ types_.FileType) *Thumb {
	t.cache.clean()
}

func (t *Thumb) GetType() types_.FileType {
	return t.typ
}

// ...
func (t *Thumb) SetSRC(src types_.URLName) *Thumb {
	t.cache.clean()
}

func (t *Thumb) GetSRC() types_.URLName {
	return types_.URLName(t.url_path)
}


// ...
func (t *Thumb) StatusManifest() bool // ( string, bool /*true - used*/ )

// ...
func (T *Thumb) GetTAG() string

// ...
func (t *Thumb) SetTypeImage( typ types_.FileType ) *Thumb {
	t.cache.clean()
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






// URLExists : проверка наличия превьюхи в запросе 
// http.Request.URL.Path -> URLpath
func URLExists( URLpath string, thumbs map[string /*nameurl*/]*Thumb ) ( *Thumb, bool /*exists*/ ) {
	t, ok := thumbs[URLpath]
	return t, ok
}

