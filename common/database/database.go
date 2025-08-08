package database

import (
    "fmt"
    bolt "go.etcd.io/bbolt"
    "github.com/labstack/gommon/log"
    "os"
    "path/filepath"
)

type Database struct {
	*bolt.DB
}

func New(name string) (*Database, error) {
	log.Debug("creating new database ", name+".db")
	// Get our directory path
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// Create our new database file
	db, err := bolt.Open(dir+"/"+name+".db", 0600, nil)
	if err != nil {
		return nil, err
	}
	return &Database{
		db,
	}, nil
}

func (d *Database) CreateBucket(name string) error {
	return d.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		log.Debug("created DB bucket ", name)
		return nil
	})

}

func (d *Database) DeleteBucket(bucket string) error {
	return d.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
}
