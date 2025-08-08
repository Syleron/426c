package main

import (
    "encoding/json"
    "errors"
    bolt "go.etcd.io/bbolt"
    "github.com/labstack/gommon/log"
    "github.com/syleron/426c/common/models"
    "time"
)

func dbUserBlockCredit(username string, total int) (int, error) {
	user, err := dbUserGet(username)
	if err != nil {
		return 0, err
	}
	user.Blocks += total
    err = db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("users"))
        // Marshal user data into bytes.
        buf, err := json.Marshal(user)
        if err != nil {
            return err
        }
        return b.Put([]byte(user.Username), buf)
    })
	if err != nil {
		log.Error(err)
	}
	return user.Blocks, err
}

func dbUserBlockDebit(username string, total int) (int, error) {
	user, err := dbUserGet(username)
	if err != nil {
		return 0, err
	}
    if (user.Blocks - total) < 0 {
		return user.Blocks, errors.New("insufficient blocks")
	} else {
		user.Blocks -= total
	}
    err = db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("users"))
        // Marshal user data into bytes.
        buf, err := json.Marshal(user)
        if err != nil {
            return err
        }
        return b.Put([]byte(user.Username), buf)
    })
	return user.Blocks, err
}

func dbUserAdd(u *models.User) error {
    t := time.Now()
    // Check existence by username key
    if _, err := dbUserGet(u.Username); err == nil {
        return errors.New("user already exists")
    }
    return db.Update(func(tx *bolt.Tx) error {
        users := tx.Bucket([]byte("users"))
        // Generate ID for the user.
        id, _ := users.NextSequence()
        u.ID = int(id)
        u.RegisteredDate = t
        // Marshal user
        buf, err := json.Marshal(u)
        if err != nil {
            return err
        }
        // Store user under username key
        if err := users.Put([]byte(u.Username), buf); err != nil {
            return err
        }
        return nil
    })
}

func dbUserGet(username string) (models.User, error) {
    var user models.User
    var ub []byte
    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("users"))
        if b == nil { return errors.New("users bucket missing") }
        v := b.Get([]byte(username))
        if v == nil { return errors.New("user does not exist") }
        ub = make([]byte, len(v))
        copy(ub, v)
        return nil
    })
    if err != nil {
        return models.User{}, err
    }
    if err := json.Unmarshal(ub, &user); err != nil {
        return models.User{}, err
    }
    return user, nil
}

func dbMessageTokenGet(uid string) (models.MessageTokenModel, error) {
    var mt models.MessageTokenModel
    var vb []byte
    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("message_tokens"))
        if b == nil { return errors.New("message_tokens bucket missing") }
        v := b.Get([]byte(uid))
        if v == nil { return errors.New("not found") }
        vb = make([]byte, len(v))
        copy(vb, v)
        return nil
    })
    if err != nil {
        return models.MessageTokenModel{}, err
    }
    if err := json.Unmarshal(vb, &mt); err != nil {
        return models.MessageTokenModel{}, err
    }
    return mt, nil
}

func dbMessageTokenAdd(mt *models.MessageTokenModel) error {
    return db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("message_tokens"))
        if b == nil { return errors.New("message_tokens bucket missing") }
        buf, err := json.Marshal(mt)
        if err != nil { return err }
        return b.Put([]byte(mt.UID), buf)
    })
}

func dbMessageTokenDelete(uid string) error {
    return db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("message_tokens"))
        if b == nil { return errors.New("message_tokens bucket missing") }
        return b.Delete([]byte(uid))
    })
}
