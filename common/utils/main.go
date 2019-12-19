package utils

import (
	"encoding/binary"
	"encoding/json"
)

// itob returns an 8-byte big endian representation of v.
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
