package checks

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 25 July 2023
 */
import (
	types_ "github.com/jhekau/favicon/internal/core/types"
)

// export cache check
func ExportCheckCache() *struct{
	Status func(fpath types_.FilePath, source_typ types_.FileType, thumb_size int) (bool, error)
	SetErr func(fpath types_.FilePath, source_typ types_.FileType, thumb_size int, err error) error
}{
	return &check_cache
}


