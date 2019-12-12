package main

import (
	"github.com/boltdb/bolt"
	"os"
	"path/filepath"
)

var db *bolt.DB

func loadDatabase() error {
	var err error
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	db, err = bolt.Open(dir+"/426c.db", 0600, nil)
	bucketCreate("users")
	if err != nil {
		return err
	}
	return nil
}