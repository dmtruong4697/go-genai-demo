package routes

import (
	"go-genai-demo/src/controllers"
	"go-genai-demo/src/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	geminiService := services.NewGeminiService()
	chatController := controllers.NewChatController(geminiService)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.POST("/chat", chatController.Chat)
		api.POST("/chat/stream", chatController.StreamChat)
		api.POST("/models", chatController.GetModels)
	}
}
