package files

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 March 2023
 */
import (
	"os"

	// err_ "github.com/jhekau/favicon/internal/core/err"
	logger_ "github.com/jhekau/favicon/internal/core/logger"
	types_ "github.com/jhekau/favicon/internal/core/types"
)

const (
	logP = `/internal/storage/files/stat.go`
	// logS01 = `S01: `
	// logS02 = `S02: `
	logS03 = `S03: os stat source image`
	logS04 = `S04: save thumb image is a folder`
)

// for test 
var (
	osStat = os.Stat
)

func Read(fpath types_.FilePath) (*os.File, error) {
	return os.Open(fpath.String())
}

func IsExists( fpath types_.FilePath, l *logger_.Logger ) ( bool, error ) {

	if f, err := osStat(fpath.String()); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, l.Typ.Error(logP, logS03, err)

	} else if f.IsDir() {
		return false, l.Typ.Error(logP, logS04)
	}

	return true, nil
}


