package routes

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/api/v1"
	// "interastral-peace.com/alnitak/internal/middleware"
)

func CollectSensitiveWordRoutes(r *gin.RouterGroup) {
	sensitiveWordGroup := r.Group("sensitive-word")

	// 敏感词检测和替换（公开接口）
	sensitiveWordGroup.POST("/check", api.CheckSensitiveWord)     // 检测敏感词
	sensitiveWordGroup.POST("/replace", api.ReplaceSensitiveWord) // 替换敏感词

	// 需要管理员权限的接口
	sensitiveWordAuth := sensitiveWordGroup.Use()
	{
		// 敏感词管理
		sensitiveWordAuth.POST("/add", api.AddSensitiveWord)         // 添加敏感词
		sensitiveWordAuth.DELETE("/delete", api.DeleteSensitiveWord) // 删除敏感词
		sensitiveWordAuth.GET("/list", api.GetSensitiveWordList)     // 获取敏感词列表
	}
}
