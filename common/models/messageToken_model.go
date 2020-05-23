package models

type MessageTokenModel struct {
	// ID of the record
	ID int `json:"id,omitempty"`
	// ID of the record
	UID string `json:"id,omitempty"`
	// Sent/received status
	Success bool `json:"success,omitempty"`
}
