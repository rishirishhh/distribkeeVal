package db

import (
	"fmt"

	bolt "github.com/boltdb/bolt"
)

// Database is an open bolt database.
type Database struct {
	db *bolt.DB
}

var defaultBucket = []byte("default")

// NewDatabase returns an instance of a database that we can work with.
func NewDatabase(dbPath string) (db *Database,closeFunc func() error,err error) {
	boltDb, err := bolt.Open(dbPath, 0600,nil)
	if err != nil {
		return nil, nil , err
	}

	db = &Database{db : boltDb}
	closeFunc = boltDb.Close

	if err := db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil , fmt.Errorf("creating default bucket: %w", err)
	}
	return db, closeFunc, nil
}

func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})

}
// SetKey sets the key to the requested value or returns an error. 
func (d *Database) SetKey(key string, value []byte) error{
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})
}

func (d *Database) GetKey(key string ) ([]byte, error){
	var result []byte
	
	err :=  d.db.View(func(tx *bolt.Tx) error {
		b:= tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		return nil
	})	

	if err == nil {
		return result , nil
	}

	return nil , err
} 

// DeleteExtraKeys deletes the keys that do not belong to this shard.
func (d *Database) DeleteExtraKeys(isExtra func(string) bool) error {
	var keys []string

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.ForEach(func(k, v []byte) error {
			ks := string(k)
			if isExtra(ks) {
				keys = append(keys, ks)
			}
			return nil
		})
	})

	if err != nil {
		return err
	}

	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)

		for _, k := range keys {
			if err := b.Delete([]byte(k)); err != nil {
				return err
			}
		}
		return nil
	})
}