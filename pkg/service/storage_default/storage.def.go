package storagedefault

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 14 August 2023
 */
import (
	files_ "github.com/jhekau/favicon/internal/storage/files"
)

// Для определения директории, в которой будут храниться нарезанные иконки
func SetFolderIcons(f string) {
	files_.SetFolderIcons(f)
}