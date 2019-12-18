package models

// Levels
// 0 - Fatal
// 1 - Error

type ErrorModel struct {
	Level int `json:"level"`
	Message int `json:"message"`
}
