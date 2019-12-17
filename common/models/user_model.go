package models

import "time"

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	PassHash       string    `json:"passHash"`
	EncPrivKey     string    `json:"encPrivKey"`
	PubKey         string    `json:"pubKey"`
	RegisteredDate time.Time `json:"registeredDate"`
	Access         int32     `json:"access"`
}
