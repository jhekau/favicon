package thumb

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 10 August 2023
 */

const (
	LogTP  = logTP
	LogT01 = logT01
	LogT03 = logT03
	LogT10 = logT10
	LogT11 = logT11
)

var TestTypEmpty Typ = ``

// замена кеша для тестирования
func (t *Thumb) TestCacheSwap( c cache ) *Thumb {
	t.cache = c
	return t
}
func (t *Thumb) GetOriginalStorageObj() StorageOBJ {
	return t.original_get().obj
}