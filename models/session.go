package models

type Session struct {
	Model    string
	Messages []Message
}

type Message struct {
	Role        Role
	Content     string
	ContentType string
}

type Role string

const (
	SYSTEM    Role = "system"
	USER      Role = "user"
	ASSISTANT Role = "assistant"
	TOOL      Role = "tool"
)
