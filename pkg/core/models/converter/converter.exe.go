package converter

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 11 August 2023
 */
import (
	"github.com/jhekau/favicon/pkg/core/types"
	"github.com/jhekau/favicon/pkg/core/models/storage"
)

type ConverterExec interface{
	Proc(source, save storage.StorageOBJ, size_px int, typ types.FileType) error
}
