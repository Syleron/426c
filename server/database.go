package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/labstack/gommon/log"
	"os"
	"path/filepath"
)

var db *bolt.DB

func loadDatabase() error {
	log.Debug("loading database 426c.db")
	var err error
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	db, err = bolt.Open(dir+"/426c.db", 0600, nil)
	bucketCreate("users")
	if err != nil {
		return err
	}
	return nil
}

func bucketCreate(name string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		log.Debug("created DB bucket ", name)
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
