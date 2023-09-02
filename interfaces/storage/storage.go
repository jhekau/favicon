package storage

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 11 August 2023
 */
import()

type Storage interface {
	NewObject(key any) (StorageOBJ, error)
}
