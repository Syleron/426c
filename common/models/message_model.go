package models

import "time"

type Message struct {
	// ID of the record
	ID int `json:"id,omitempty"`
	// Recipients version of encrypted message
	ToMessage string `json:"toMessage,omitempty"`
	// Senders version of encrypted message
	FromMessage string `json:"fromMessage,omitempty"`
	// Sent from username
	From string `json:"from,omitempty"`
	// Sent to username
	To string `json:"to,omitempty"`
	// Sent time/date
	Date time.Time `json:"date,omitempty"`
	// Sent/received status
	Success bool `json:"success,omitempty"`
}

type MsgToRequestModel struct {
	Message `json:"message,omitempty"`
}

type MsgToResponseModel struct {
	Success bool `json:"success,omitempty"`
	MsgID   int  `json:"msgID,omitempty"`
	To string `json:"to,omitempty"`
}

type MsgResponseModel struct {
	Message `json:"message,omitempty"`
}
