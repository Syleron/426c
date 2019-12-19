package models

type RegisterRequestModel struct {
	Username   string `json:"username"`
	PassHash   string `json:"passhash"`
	EncPrivKey string `json:"encprivkey"`
	PubKey     string `json:"pubkey"`
}

type RegisterResponseModel struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}