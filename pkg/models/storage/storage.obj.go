package storage

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 11 August 2023
 */
import(
	"io"
)

type StorageOBJ interface{
	Reader() (io.ReadSeekCloser , error)
	Writer() (io.WriteCloser, error)
	Key() StorageKey
	IsExists() ( bool, error )
}
