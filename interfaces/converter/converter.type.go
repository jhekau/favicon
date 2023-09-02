package converter

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 11 August 2023
 */
import (
	"github.com/jhekau/favicon/interfaces/storage"
	"github.com/jhekau/favicon/domain/types"
)

type ConverterTyp interface{
	Do(source, save storage.StorageOBJ, size_px int, typThumb types.FileType) (complete bool, err error)
}
