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

//func bucketGet(bucket string) ([]Issue, error) {
//	response := []Issue{}
//	err := db.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(bucket))
//		if b == nil {
//			return errors.New("unable to fetch bucket")
//		}
//		c := b.Cursor()
//		for k, v := c.First(); k != nil; k, v = c.Next() {
//			var issue Issue
//			if err := json.Unmarshal(v, &issue); err != nil {
//				return err
//			}
//			if !issue.Resolved {
//				if issue.History == nil {
//					issue.History = make([]History, 0)
//				}
//				response = append(response, issue)
//			}
//		}
//		return nil
//	})
//	return response, err
//}