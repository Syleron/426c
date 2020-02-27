package models

// TC/S
// Registered
// Registration status (open/closed)

type StatusResponseModel struct {
	MsgCost int    `json:"msgCost,omitempty"`
	OnlineUsers int `json:"onlineUsers,omitempty"`
}
