package favicon

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 09 March 2023
 */
import (
	"sync"

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
	original_img types_.FilePath // jpg, png ... 
	folder_work types_.Folder
	
}

