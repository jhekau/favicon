package favicon

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 09 March 2023
 */
import (
	"net/http"
	"sync"

	thumb_ "github.com/jhekau/favicon/thumb"
	types_ "github.com/jhekau/favicon/types"
)

var (
	// ~~ interface ~~ 

	DefaultList = func(){}
	CustomList = func(){}
)

type Thumbs struct {
	sync.RWMutex
	original_svg types_.FilePath
	original_img types_.FilePath
	folder_work types_.Folder
	thumbs map[types_.URLHref/*clear*/]*thumb_.Thumb
}

func (t *Thumbs) Append(*thumb_.Thumb) *Thumbs
func (t *Thumbs) SetFolderWork( folder string ) *Thumbs
func (t *Thumbs) SetFilepathSourceSVG( fpath string ) *Thumbs
func (t *Thumbs) SetFilepathSourceIMG( fpath string ) *Thumbs

func (t *Thumbs) Handle( w http.Response, r *http.Request )
func (t *Thumbs) ServeFile( URLpath string ) ( fpath string, exists bool, err error )

func (t *Thumbs) server_file_thumb( URLpath string )
func (t *Thumbs) server_file_manifest( URLpath string )
func (T *Thumbs) manifest_url_exists( URLpath string ) ( fpath string, exists bool, err error )
func (t *Thumbs) get_thumbs() map[types_.URLHref/*clear*/]*thumb_.Thumb

