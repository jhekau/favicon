package favicon

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 09 March 2023
 */
import (
	"net/http"
	"net/url"
	"sync"

	manifest_ "github.com/jhekau/favicon/manifest"
	thumb_ "github.com/jhekau/favicon/thumb"
	types_ "github.com/jhekau/favicon/types"
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
	t.thumbs[thumb.GetHREFClear()] = thumb
	return t
}

func (t *Thumbs) set_folder_work( folder string ) *Thumbs
func (t *Thumbs) get_folder_work() types_.Folder
func (t *Thumbs) set_filepath_source_svg( fpath string ) *Thumbs
func (t *Thumbs) get_filepath_source_svg() types_.FilePath {
	return t.source_svg
}
func (t *Thumbs) set_filepath_source_img( fpath string ) *Thumbs
func (t *Thumbs) get_filepath_source_img() types_.FilePath {
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
			// error
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
		// return error
	} else if exists {
		return manifest.String(), true, nil
	}

	if thumb, exists, err := t.server_file_thumb(url_); err != nil {
		// return error
	} else if exists {
		return thumb.String(), true, nil
	}

	return ``, false, nil
}

func (t *Thumbs) server_file_thumb( url_ *url.URL ) (fpath types_.FilePath, exists bool, err error) {

	thumb, exists := thumb_.URLExists(url_, t.thumbs)
	if !exists {
		return ``, false, nil
	}

	fpath, err = thumb.GetFile(t.get_folder_work(), t.get_filepath_source_img(), t.get_filepath_source_svg())
	if err != nil {
		// return error
	}

	return fpath, true, nil
}

func (t *Thumbs) server_file_manifest( url_ *url.URL ) (manifest types_.FilePath, exists bool, err error) {

	if !t.manifest.URLExists(url_) {
		return ``, false, nil
	}

	manifest, exists, err = t.manifest.GetFile(t.get_folder_work(), t.thumbs)
	if err != nil {
		// return error
	}
	return manifest, exists, nil
}

// func (T *Thumbs) manifest_url_exists( URLpath string ) ( fpath string, exists bool, err error )
func (t *Thumbs) get_thumbs() map[types_.URLHref/*clear*/]*thumb_.Thumb

func (t *Thumbs) tags_html() string


func default_list( ) *Thumbs
func custom_list( thumbs ...*Thumbs ) *Thumbs
