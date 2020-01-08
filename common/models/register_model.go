package models

type RegisterRequestModel struct {
	Username   string `json:"username,omitempty"`
	PassHash   string `json:"passhash,omitempty"`
	EncPrivKey string `json:"encprivkey,omitempty"`
	PubKey     string `json:"pubkey,omitempty"`
}

type RegisterResponseModel struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}
