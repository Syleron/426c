package models

import "time"

type Message struct {
	ID      int       `json:"id,omitempty"`
	Message string    `json:"message,omitempty"`
	From    string    `json:"from,omitempty"`
	To      string    `json:"to,omitempty"`
	Date    time.Time `json:"date,omitempty"`
	Success bool      `json:"success,omitempty"`
}

type MsgToRequestModel struct {
	Message Message `json:"message,omitempty"`
}

type MsgToResponseModel struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}
