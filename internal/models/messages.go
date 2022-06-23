package models

type Message struct {
	Message string `json:"message"`
}

type MessageDB struct {
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
