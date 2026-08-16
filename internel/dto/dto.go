package dto

type Message struct {
	Type    string `json:"type"`
	Author  string `json:"author"`
	Payload string `json:"payload"`
}
