package main

import (
	"fmt"
	"github.com/boltdb/bolt"
)

func bucketCreate(name string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func bucketDelete(bucket string) error {
	return db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
}
