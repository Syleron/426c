package models

type RegisterModel struct {
	Username   string `json:"username"`
	PassHash   string `json:"passhash"`
	EncPrivKey string `json:"encprivkey"`
	PubKey     string `json:"pubkey"`
}
