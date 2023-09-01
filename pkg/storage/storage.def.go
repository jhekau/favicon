package storagedefault

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 August 2023
 */
import (
	files_ "github.com/jhekau/favicon/internal/storage/files"
)

// Для переопределния дефолтовой директории, в которой будут храниться нарезанные иконки для интегрированного storage
func SetFolderIcons(f string) {
	files_.DirIconsDefault = f
}