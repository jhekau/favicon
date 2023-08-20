package checks

/* *
 * Copyright (c) 2023, @jhekau <mr.evgeny.u@gmail.com>
 * 12 August 2023
 */
import (
	"fmt"
	"sync"
	storage_ "github.com/jhekau/favicon/pkg/core/models/storage"
)

type CacheStatus struct {
	m sync.Map
}
func (c *CacheStatus) cache_key(original storage_.StorageKey, originalSVG bool, thumb_size int) string {
	return fmt.Sprintf(`%s, %v, %d`, original, originalSVG, thumb_size)
}
func (c *CacheStatus) Status(original storage_.StorageKey, originalSVG bool, thumb_size int) (bool, error) {

	e, ok := c.m.Load(c.cache_key(original, originalSVG, thumb_size))

	var err error
	if e != nil {
		err = e.(error)
	}

	return ok, err
}
func (c *CacheStatus) SetErr(original storage_.StorageKey, originalSVG bool, thumb_size int, err error) error {
	c.m.Store(c.cache_key(original, originalSVG, thumb_size), err)
	return err
}