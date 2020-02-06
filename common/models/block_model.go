package models

type BlockResponseModel struct {
	Blocks  int    `json:"blocks,omitempty"`
	MsgCost int    `json:"msgCost,omitempty"`
}
