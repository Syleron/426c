package models

import "time"

type MsgToRequestModel struct {
	Message string    `json:"message,omitempty"`
	From    string    `json:"from,omitempty"`
	To      string    `json:"to,omitempty"`
	Date    time.Time `json:"date,omitempty"`
}

type MsgToResponseModel struct {
}
