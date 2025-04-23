package models

type ChatRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}