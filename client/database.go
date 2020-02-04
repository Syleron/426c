package main

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/syleron/426c/common/models"
	"github.com/syleron/426c/common/utils"
)

// dbMessageAdd - Add a message to our data store
func dbMessageAdd(m *models.Message) (int, error) {
	var msgID int
	if db == nil {
		return 0, nil
	}
	// make sure our bucket exists before attempting to add a message
	db.CreateBucket(m.To)
	// Attempt to add our message
	return msgID, db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		b := tx.Bucket([]byte(m.To))
		// Generate ID for the user.
		id, _ := b.NextSequence()
		// Set our ID
		m.ID = int(id)
		// Set our local var
		msgID = m.ID
		// Marshal user data into bytes.
		buf, err := json.Marshal(m)
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put(utils.Itob(m.ID), buf)
	})
}

func dbMessagesGet(toUsername string, fromUsername string) ([]models.Message, error) {
	if db == nil {
		return nil, nil
	}
	var messages []models.Message

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(toUsername))

		if b == nil {
			return errors.New("bucket doesn't exist")
		}

		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var found models.Message

			// only return 50 messages
			if len(messages) >= 50 {
				return nil
			}

			// copy data into our issue object
			if err := json.Unmarshal(v, &found); err != nil {
				return err
			}

			// Make sure our from and to match
			if found.To == toUsername && found.From == fromUsername ||
				found.To == fromUsername && found.From == toUsername {
				// Add our message to the array
				messages = append(messages, found)
			}
		}
		return nil
	})

	// return our issue
	return messages, err
}

func dbUserList() ([]models.User, error) {
	if db == nil {
		return []models.User{}, nil
	}
	response := []models.User{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b == nil {
			return errors.New("unable to fetch bucket")
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var user models.User
			if err := json.Unmarshal(v, &user); err != nil {
				return err
			}
			response = append(response, user)
		}
		return nil
	})
	return response, err
}

func dbUserAdd(u models.User) error {
	if db == nil {
		return nil
	}
	// make sure our bucket exists before attempting to add a message
	db.CreateBucket("users")
	// Check to see if our user exists
	_, err := dbUserGet(u.Username)
	if err == nil {
		return errors.New("user already exists")
	}
	// Attempt to add our message
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		b := tx.Bucket([]byte("users"))
		// Generate ID for the user.
		id, _ := b.NextSequence()
		// Set our ID
		u.ID = int(id)
		// Marshal user data into bytes.
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put(utils.Itob(u.ID), buf)
	})
}

func dbUserGet(username string) (models.User, error) {
	if db == nil {
		return models.User{}, nil
	}
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
