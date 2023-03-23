package manifest

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 March 2023
 */
import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	thumb_ "github.com/jhekau/favicon/thumb"
	types_ "github.com/jhekau/favicon/types"
)

var (
	// ~~ interface ~~

	Default = default_get
)

const (
	default_name_url types_.URLHref = `manifest.webmanifest`
	default_name_file types_.FileName = `manifest.webmanifest`
)

type Manifest struct {
	sync.RWMutex
	cache cache

	url_href types_.URLHref
	url_href_clear types_.URLHref
	filename string
}

func (m *Manifest) SetNameURL(src string) *Manifest {
	return m.set_name_url(src)
}

func (m *Manifest) GetNameURL() types_.URLHref {
	return m.get_name_url()
}

func (m *Manifest) GetNameURLClear() types_.URLHref {
	return m.get_name_url_clear()
}

func (m *Manifest) SetNameFile(f string) *Manifest {
	return m.set_name_file(f)
}

func (m *Manifest) GetNameFile() types_.FileName {
	return m.get_name_file()
}

func (m *Manifest) URLExists(url_ *url.URL) bool {
	return m.url_exists(url_)
}

func (m *Manifest) GetFile(
	folder_work types_.Folder,
	thumbs map[types_.URLHref]*thumb_.Thumb,
)(
	fpath types_.FilePath,
	exists bool,
	err error,
){
	return m.get_file(folder_work, thumbs)
}





//
func (m *Manifest) set_name_url(src string) *Manifest {

	// if m == nil {
	// 	m = &Manifest{}
	// }
	m.Lock()
	defer m.Unlock()
	m.cache.clean()
	
	m.url_href = types_.URLHref(src)
	{
		u, err := url.Parse(`http://domain.com`)
		if err != nil {
			// error
		} else {
			m.url_href_clear = types_.URLHref(u.JoinPath(src).Path)
		}
	}

	return m
}
func (m *Manifest) get_name_url() types_.URLHref {
	
	m.RLock()
	defer m.RUnlock()
	
	return m.url_href
}

func (m *Manifest) get_name_url_clear() types_.URLHref {
	
	m.RLock()
	defer m.RUnlock()

	return m.url_href_clear
}

//
func (m *Manifest) set_name_file(f string) *Manifest {
	
	m.Lock()
	defer m.Unlock()

	m.cache.clean()

	m.filename = f
	return m
}

func (m *Manifest) get_name_file() types_.FileName {

	m.RLock()
	defer m.RUnlock()
	
	return types_.FileName(m.filename)
}

func (m *Manifest) generate(
	thumbs map[types_.URLHref]*thumb_.Thumb,
)(
	filebody []byte,
	status_generate bool,
	err error,
){

	m.RLock()
	defer m.RUnlock()

	if len(thumbs) == 0 {
		return nil, false, nil
	}

	/* ------------------------- filebody -------------------------------- /*
	{
		"icons": [
			{ "src": "/icon-192.png", "type": "image/png", "sizes": "192x192" },
			{ "src": "/icon-512.png", "type": "image/png", "sizes": "512x512" }
		]
	}
	/* -------------------------------------------------------------------- */
	type icon struct {
		Src string `json:"src"`
		Type string `json:"type"`
		Size int `json:"sizes"`
	}

	type icons struct {
		Icons []icon `json:"icons"`
	}

	list := make([]icon, 0)
	for _, thumb := range thumbs {
		if thumb.StatusManifest() {
			list = append(list, icon{
				Src: thumb.GetHREF().String(),
				Type: thumb.GetType().String(),
				Size: int(thumb.GetSize()),
			})
		}
	}

	if len(list) == 0 {
		return nil, false, nil
	}

	body, err := json.Marshal(list)
	if err != nil {
		// return error
	}
	return body, true, nil
}


//
func (m *Manifest) file_create(
	folder_work types_.Folder,
	thumbs map[types_.URLHref]*thumb_.Thumb,
)(
	fpath types_.FilePath,
	status_create bool,
	err error,
){
	filebody, status, err := m.generate(thumbs)
	if err != nil {
		// return error
	}
	if !status {
		return ``, false, nil
	}

	filename := ``
	m.RLock()
	filename = m.filename
	m.RUnlock()
	
	if filename == `` {
		m.Lock()
		m.filename = strconv.FormatInt(time.Now().Unix(), 10)+`.manifest`
		m.Unlock()
	}

	fpath = types_.FilePath(
		filepath.Join(folder_work.String(), m.filename),
	)
	os.Remove(fpath.String())

	if err = ioutil.WriteFile(fpath.String(), filebody, 0775); err != nil {
		// return ``, false, error
	}
	return fpath, true, nil
}


//
func (m *Manifest) get_file(
	folder_work types_.Folder,
	thumbs map[types_.URLHref]*thumb_.Thumb,
)(
	fpath types_.FilePath,
	exists bool,
	err error,
){
	already_exists := m.cache.get_file_exists()
	switch already_exists {
	case types_.FileExistsOK:
		if fpath := m.cache.get_filepath(); fpath != `` {
			return fpath, true, nil
		}
	case types_.FileExistsNOT:
		return ``, false, nil
	}

	fpath, state_create, err := m.file_create(folder_work, thumbs)
	if err != nil {
		// return error
	}
	if !state_create {
		m.cache.set_file_exists(types_.FileExistsNOT)
		return ``, false, nil
	}
	m.cache.set_filepath(fpath)
	m.cache.set_file_exists(types_.FileExistsOK)
	return fpath, true, nil
}



// 
func (m *Manifest) url_exists(url_ *url.URL) bool {

	// if m == nil {
	// 	return false
	// }

	/*
	src := url_.Path
	manifest := m.GetNameURLClear().String()
	if src == manifest {
		return true
	} else if src[0] != '/' && `/`+src == manifest {
		return true
	}
	*/
	return url_.Path == m.GetNameURLClear().String()
}


func default_get() *Manifest {
	return (&Manifest{}).
				SetNameURL(default_name_url.String()).
				SetNameFile(default_name_file.String())
}


//
type cache struct {
	sync.RWMutex
	filepath types_.FilePath
	file_exists types_.FileExists
}

func (c *cache) clean() {
	c = &cache{}
}

//
func (c *cache) set_file_exists( state types_.FileExists ) {
	c.file_exists = state
}

func (c *cache) get_file_exists() types_.FileExists {
	return c.file_exists
}

//
func (c *cache) set_filepath( fpath types_.FilePath ) {
	c.filepath  = fpath
}

func (c *cache) get_filepath() types_.FilePath {
	return c.filepath
}
