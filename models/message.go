package models

type Message struct {
	UserID      int    `json:"userid"`
	UserMessage string `json:"usermessage"`
	Response    string `json:"response"`
	Timestamp   string `json:"timestamp"`
	Path        string `json:"path"`
}
