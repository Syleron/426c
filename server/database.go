package main

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/syleron/426c/common/models"
	"github.com/syleron/426c/common/utils"
	"time"
)

func dbUserBlockDebit() {}

func dbUserCreditDebit() {}

func dbUserAdd(u *models.User) error {
	t := time.Now()
	_, err := dbUserGet(u.Username)
	if err == nil {
		return errors.New("user already exists")
	}
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		b := tx.Bucket([]byte("users"))
		// Generate ID for the user.
		id, _ := b.NextSequence()
		// Set our ID
		u.ID = int(id)
		// Set our reg. date
		u.RegisteredDate = t
		// Marshal user data into bytes.
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put(utils.Itob(u.ID), buf)
	})
}

func dbUserDelete() {}

func dbUserGet(username string) (models.User, error) {
	var user models.User
	var ub []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		return b.ForEach(func(k, v []byte) error {
			var found models.User

			// copy data into our issue object
			if err := json.Unmarshal(v, &found); err != nil {
				return err
			}

			if found.Username != username {
				return nil
			}

			// initiate our object
			ub = make([]byte, len(v))

			// copy our data to the object
			copy(ub, v)

			return nil
		})
	})
	// Make sure we have something
	if err != nil || len(ub) == 0 {
		return models.User{}, errors.New("user does not exist")
	}
	// unmarshal our data
	if err := json.Unmarshal(ub, &user); err != nil {
		return models.User{}, err
	}
	// return our issue
	return user, err
}

