package converter

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 11 August 2023
 */
import (
	"github.com/jhekau/favicon/pkg/models/storage"
	"github.com/jhekau/favicon/pkg/core/types"
)

type Converter interface{
	Do(source, save storage.StorageOBJ, originalSVG bool, typThumb types.FileType, size_px int) error
}
