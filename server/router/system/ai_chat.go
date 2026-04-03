package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AIChatRouter struct{}

func (r *AIChatRouter) InitAIChatRouter(Router *gin.RouterGroup) {
	aiRouter := Router.Group("ai").Use(middleware.JWTAuth(), middleware.OperationRecord())
	aiStreamRouter := Router.Group("ai").Use(middleware.JWTAuth())
	{
		aiRouter.POST("chat", aiChatApi.Chat)
		aiStreamRouter.POST("chat/stream", aiChatApi.ChatStream)
	}
}
