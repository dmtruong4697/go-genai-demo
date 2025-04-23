package controllers

import (
	"go-genai-demo/src/models"
	"go-genai-demo/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	Service *services.GeminiService
}

func NewChatController(s *services.GeminiService) *ChatController {
	return &ChatController{Service: s}
}

// Chat godoc
// @Summary chat one turn
// @Description simple question (one-turn)
// @Tags chat
// @Accept  json
// @Produce  json
// @Param request body models.ChatRequest true "prompt"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chat [post]
func (c *ChatController) Chat(ctx *gin.Context) {
	var req models.ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := c.Service.Chat(req.Prompt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": response})
}

// StreamChat godoc
// @Summary Chat streaming
// @Description stream phản hồi
// @Tags chat
// @Accept json
// @Produce text/event-stream
// @Param request body models.ChatRequest true "prompt gửi tới gemini"
// @Success 200 {string} string "Streaming response"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /chat/stream [post]
func (cc *ChatController) StreamChat(ctx *gin.Context) {
	var req models.ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := cc.Service.StreamChat(req.Prompt, ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (mc *ChatController) GetModels(c *gin.Context) {
	models, err := mc.Service.ListModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models)
}