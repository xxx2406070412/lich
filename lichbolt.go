// boltdboperation
package main

import (
	"github.com/boltdb/bolt"
	"github.com/coocood/freecache"
)

var db *bolt.DB

func init() {

	var err error
	db, err = bolt.Open(g_dataPath+"my.db", 0600, nil)
	if err != nil {
		logger.Infof(err.Error())
	}

}

func lich_bolt_insert(bBucket []byte, bkey []byte, bvalue []byte) error {

	db.Update(func(tx *bolt.Tx) error {
		// ...read or write...
		b, err := tx.CreateBucketIfNotExists(bBucket)
		if err != nil {
			return err
		}
		b.Put(bkey, bvalue)

		return nil
	})

	return nil
}

func lich_bolt_get(bBucket []byte, bkey []byte) (value []byte, err error) {

	var result []byte
	db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(bBucket)

		if nil == b {
			result = nil
		} else {
			result = b.Get(bkey)
		}

		return nil
	})

	if nil == result {
		return nil, freecache.ErrNotFound
	} else {
		return result, nil
	}
}

func lich_bolt_delete(bBucket []byte, bkey []byte) (err error) {

	db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bBucket)

		if nil != b {
			b.Delete(bkey)
		}

		return nil
	})

	return nil
}
