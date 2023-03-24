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

	defaults_ "github.com/jhekau/favicon/defaults"
	err_ "github.com/jhekau/favicon/err"
	manifest_ "github.com/jhekau/favicon/manifest"
	thumb_ "github.com/jhekau/favicon/thumb"
	types_ "github.com/jhekau/favicon/types"
)

const (
	logT01 = `T01: get serve file`
	logT02 = `T02: get file manifest`
	logT03 = `T03: get file thumb`
	logT04 = `T04: get file thumb`
	logT05 = `T05: get file manifest`
)
func errT(i... interface{}) error {
	return err_.Err(err_.TypeError, `/thumbs.go`, i)
} 

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

type Thumbs struct {
	sync.RWMutex
	source_svg types_.FilePath
	source_img types_.FilePath
	folder_work types_.Folder
	thumbs map[types_.URLHref/*clear*/]*thumb_.Thumb
	manifest manifest_.Manifest
}

func (t *Thumbs) Append(thumb *thumb_.Thumb) *Thumbs {
	return t.append(thumb)
}
func (t *Thumbs) SetFolderWork( folder string ) *Thumbs {
	return t.set_folder_work(folder)
}
func (t *Thumbs) SetFilepathSourceSVG( fpath string ) *Thumbs {
	return t.set_filepath_source_svg( fpath )
}
func (t *Thumbs) SetFilepathSourceIMG( fpath string ) *Thumbs {
	return t.set_filepath_source_img( fpath )
}
func (t *Thumbs) Handle() {
	t.handle()
}
func (t *Thumbs) ServeFile( url_ *url.URL ) ( fpath string, exists bool, err error ) {
	return t.serve_file(url_)
}
func (t *Thumbs) TagsHTML() string {
	return t.tags_html()
}


func (t *Thumbs) append(thumb *thumb_.Thumb) *Thumbs {
	t.Lock()
	defer t.Unlock()
	
	t.thumbs[thumb.GetHREFClear()] = thumb
	return t
}

func (t *Thumbs) set_folder_work( folder string ) *Thumbs {
	t.Lock()
	defer t.Unlock()

	t.folder_work = types_.Folder(folder)
	return t
}
func (t *Thumbs) get_folder_work() types_.Folder {
	t.RLock()
	defer t.RUnlock()

	return t.folder_work
}
func (t *Thumbs) set_filepath_source_svg( fpath string ) *Thumbs {
	t.Lock()
	defer t.RUnlock()

	t.source_svg = types_.FilePath(fpath)
	return t
}
func (t *Thumbs) get_filepath_source_svg() types_.FilePath {
	t.RLock()
	defer t.RUnlock()

	return t.source_svg
}
func (t *Thumbs) set_filepath_source_img( fpath string ) *Thumbs {
	t.Lock()
	defer t.Unlock()

	t.source_img = types_.FilePath(fpath)
	return t
}
func (t *Thumbs) get_filepath_source_img() types_.FilePath {
	t.RLock()
	defer t.RUnlock()

	return t.source_img
}

func (t *Thumbs) handle() {
	http.Handle(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		switch r.Method {
		case http.MethodGet, http.MethodPost:
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		fpath, exists, err := t.ServeFile(r.URL)
		if err != nil {
			log.Println(errT(logT01, err))
			w.WriteHeader(http.StatusInternalServerError)
		} else if !exists {
			w.WriteHeader(http.StatusNotFound)
		} else {
			http.ServeFile(w, r, fpath)
		}
	}))
}

func (t *Thumbs) serve_file( url_ *url.URL ) ( fpath string, exists bool, err error ) {

	if manifest, exists, err := t.server_file_manifest(url_); err != nil {
		return ``, false, errT(logT02, err)
	} else if exists {
		return manifest.String(), true, nil
	}

	if thumb, exists, err := t.server_file_thumb(url_); err != nil {
		return ``, false, errT(logT03, err)
	} else if exists {
		return thumb.String(), true, nil
	}

	return ``, false, nil
}

func (t *Thumbs) server_file_thumb( url_ *url.URL ) (fpath types_.FilePath, exists bool, err error) {

	t.RLock()
	thumb, exists := thumb_.URLExists(url_, t.thumbs)
	t.RUnlock()

	if !exists {
		return ``, false, nil
	}

	fpath, err = thumb.GetFile(t.get_folder_work(), t.get_filepath_source_img(), t.get_filepath_source_svg())
	if err != nil {
		return ``, false, errT(logT04, err)
	}

	return fpath, true, nil
}

func (t *Thumbs) server_file_manifest( url_ *url.URL ) (manifest types_.FilePath, exists bool, err error) {

	if !t.manifest.URLExists(url_) {
		return ``, false, nil
	}

	manifest, exists, err = t.manifest.GetFile(t.get_folder_work(), t.thumbs)
	if err != nil {
		return ``, false, errT(logT05, err)
	}
	return manifest, exists, nil
}

// func (T *Thumbs) manifest_url_exists( URLpath string ) ( fpath string, exists bool, err error )
// func (t *Thumbs) get_thumbs() map[types_.URLHref/*clear*/]*thumb_.Thumb

func (t *Thumbs) tags_html() string {

	var tags strings.Builder

	t.RLock()
	for _, thumb := range t.thumbs {
		tags.WriteString(thumb.GetTAG())
	}
	tags.WriteString(t.manifest.GetTAG())
	t.RUnlock()

	return tags.String()
}


func default_list() *Thumbs {
	return &Thumbs{
		thumbs: func() map[types_.URLHref]*thumb_.Thumb {
			m := map[types_.URLHref]*thumb_.Thumb{}
			for _, thumb := range defaults_.Defaults() {
				m[thumb.GetHREFClear()] = thumb
			}
			return m
		}(),
	}
}

func custom_list( thumbs ...*thumb_.Thumb ) *Thumbs {
	return &Thumbs{
		thumbs: func() map[types_.URLHref]*thumb_.Thumb {
			m := map[types_.URLHref]*thumb_.Thumb{}
			for _, thumb := range thumbs {
				m[thumb.GetHREFClear()] = thumb
			}
			return m
		}(),
	}
}
