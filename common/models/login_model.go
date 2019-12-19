package models

type LoginRequestModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseModel struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}
