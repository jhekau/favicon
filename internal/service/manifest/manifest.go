package manifest

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 March 2023
 */

type Manifest struct {}

/*
import (
	"encoding/json"
	// "io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	logger_ "github.com/jhekau/favicon/pkg/core/models/logger"
	thumb_ "github.com/jhekau/favicon/internal/service/thumb"
	types_ "github.com/jhekau/favicon/internal/core/types"
)

const (
	logM   = `/manifest/manifest.go`
	logM01 = `M01: url: error parsing standart template domain.com`
	logM02 = `M02: json: marshal manifest`
	logM03 = `M03: generate manifest body`
	logM04 = `M04: write file manifest`
	logM05 = `M05: manifest file create`
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
	s sync.RWMutex
	l logger_.Logger
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

func (m *Manifest) GetTAG() string {
	return m.get_tag()
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

	m.s.Lock()
	defer m.s.Unlock()
	m.cache.clean()
	
	m.url_href = types_.URLHref(src)
	{
		u, err := url.Parse(`http://domain.com`)
		if err != nil {
			log.Println(m.l.Error(logM, logM01, err))
		} else {
			m.url_href_clear = types_.URLHref(u.JoinPath(src).Path)
		}
	}

	return m
}
func (m *Manifest) get_name_url() types_.URLHref {
	
	m.s.RLock()
	defer m.s.RUnlock()
	
	return m.url_href
}

func (m *Manifest) get_name_url_clear() types_.URLHref {
	
	m.s.RLock()
	defer m.s.RUnlock()

	return m.url_href_clear
}

//
func (m *Manifest) set_name_file(f string) *Manifest {
	
	m.s.Lock()
	defer m.s.Unlock()

	m.cache.clean()

	m.filename = f
	return m
}

func (m *Manifest) get_name_file() types_.FileName {

	m.s.RLock()
	defer m.s.RUnlock()
	
	return types_.FileName(m.filename)
}

func (m *Manifest) generate(
	thumbs map[types_.URLHref]*thumb_.Thumb,
)(
	filebody []byte,
	status_generate bool,
	err error,
){

	m.s.RLock()
	defer m.s.RUnlock()

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
	/* -------------------------------------------------------------------- * /
	type icon struct {
		Src string `json:"src"`
		Type string `json:"type"`
		Size int `json:"sizes"`
	}

	list := struct {
		Icons []icon `json:"icons"`
	}{
		Icons: make([]icon, 0),
	}

	// list := make([]icon, 0)
	for _, thumb := range thumbs {
		if thumb.StatusManifest() {
			list.Icons = append(list.Icons, icon{
				Src: thumb.GetHREF().String(),
				Type: thumb.GetType().String(),
				Size: int(thumb.GetSize()),
			})
		}
	}

	if len(list.Icons) == 0 {
		return nil, false, nil
	}

	body, err := json.Marshal(list)
	if err != nil {
		return nil, false, m.l.Error(logM, logM02, err)
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
		return ``, false, m.l.Error(logM, logM03, err)
	}
	if !status {
		return ``, false, nil
	}

	filename := ``
	m.s.RLock()
	filename = m.filename
	m.s.RUnlock()
	
	if filename == `` {
		m.s.Lock()
		m.filename = strconv.FormatInt(time.Now().Unix(), 10)+`.manifest`
		m.s.Unlock()
	}

	fpath = types_.FilePath(
		filepath.Join(folder_work.String(), m.filename),
	)
	os.Remove(fpath.String())

	if err = os.WriteFile(fpath.String(), filebody, 0775); err != nil {
		return ``, false, m.l.Error(logM, logM04, err)
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
		return ``, false, m.l.Error(logM, logM05, err)
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
func (m *Manifest) get_tag() string {
	m.s.RLock()
	defer m.s.RUnlock()

	return `<link rel="manifest" href="`+m.get_name_url().String()+`">`
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
	* /
	return url_.Path == m.GetNameURLClear().String()
}


func default_get() Manifest {
	return *(&Manifest{}).
				SetNameURL(default_name_url.String()).
				SetNameFile(default_name_file.String())
}


//
type cache struct {
	s sync.RWMutex
	filepath types_.FilePath
	file_exists types_.FileExists
}

func (c *cache) clean() {
	c = &cache{}
}

//
func (c *cache) set_file_exists( state types_.FileExists ) {
	c.s.Lock()
	c.file_exists = state
	c.s.Unlock()
}

func (c *cache) get_file_exists() types_.FileExists {
	c.s.RLock()
	file_exists := c.file_exists
	c.s.RUnlock()
	return file_exists
}

//
func (c *cache) set_filepath( fpath types_.FilePath ) {
	c.s.Lock()
	c.filepath  = fpath
	c.s.Unlock()
}

func (c *cache) get_filepath() types_.FilePath {
	c.s.RLock()
	filepath := c.filepath
	c.s.RUnlock()
	return filepath
}
*/

