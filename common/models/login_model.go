package models

type LoginRequestModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Version  string    `json:"version"`
}

type LoginResponseModel struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Blocks  int    `json:"blocks"`
	MsgCost int    `json:"msgCost"`
}
