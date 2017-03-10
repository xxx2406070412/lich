// cache
package main

import (
	"github.com/coocood/freecache"
)

var (
	g_pstcache *freecache.Cache
)

func init() {
	return
}

func lich_cache_insert(bkey []byte, bvalue []byte, iexpireSeconds int) error {

	if 0 == iexpireSeconds {
		iexpireSeconds = 60
	}

	err := g_pstcache.Set(bkey, bvalue, iexpireSeconds)

	if nil != err {
		return err
	}

	return nil
}

func lich_cache_get(bkey []byte) (value []byte, err error) {

	result, err := g_pstcache.Get(bkey)

	if err != nil {
		return nil, err
	} else {
		return result, nil
	}

}

func lich_cache_delete(bkey []byte) error {

	g_pstcache.Del(bkey)

	return nil
}
