package utils

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io/ioutil"
)

// utils.Itob returns an 8-byte big endian representation of v.
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Convert our interface to a json byte array
func MarshalResponse(obj interface{}) []byte{
	// Convert our object to a json byte array to send
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return b
}

func LoadFile(file string) ([]byte, error) {
	c, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, errors.New("unable to load file: " + file)
	}

	return []byte(c), nil
}

func WriteFile(contents string, file string) error {
	err := ioutil.WriteFile(file, []byte(contents), 0644)
	if err != nil {
		return err
	}
	return nil
}
