package models

import "time"

type User struct {
	ID             int       `json:"id,omitempty"`
	Username       string    `json:"username,omitempty"`
	PassHash       string    `json:"passHash,omitempty"`
	EncPrivKey     string    `json:"encPrivKey,omitempty"`
	PubKey         string    `json:"pubKey,omitempty"`
	RegisteredDate time.Time `json:"registeredDate,omitempty"`
	Access         int32     `json:"access,omitempty"`
	Blocks         int       `json:"blocks,omitempty"`
}

type UserRequestModel struct {
	Username string `json:"username,omitempty"`
}

type UserResponseModel struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	User    User   `json:"user,omitempty"`
}
