package main

import (
	"log"

	"notify-template/internal/api"
	"notify-template/internal/wire"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化依赖
	notifyService, err := wire.InitializeAPI()
	if err != nil {
		log.Fatal("Failed to initialize dependencies: ", err)
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// 设置路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 添加通知路由
	r.POST("/notify", func(c *gin.Context) {
		var req struct {
			Message string `json:"message"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := notifyService.SendNotification(req.Message); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": "success"})
	})

	// 模板相关路由
	templateHandler := api.NewTemplateHandler(notifyService.GetTemplateService())
	templateGroup := r.Group("/api/templates")
	{
		templateGroup.POST("", templateHandler.CreateTemplate)
		templateGroup.PUT("/:id", templateHandler.UpdateTemplate)
		templateGroup.DELETE("/:id", templateHandler.DeleteTemplate)
		templateGroup.GET("/:id", templateHandler.GetTemplate)
		templateGroup.GET("", templateHandler.ListTemplates)

		// 版本相关路由
		templateGroup.POST("/:template_id/versions", templateHandler.CreateTemplateVersion)
		templateGroup.PUT("/versions/:id", templateHandler.UpdateTemplateVersion)
		templateGroup.GET("/versions/:id", templateHandler.GetTemplateVersion)
		templateGroup.GET("/:template_id/versions", templateHandler.ListTemplateVersions)
		templateGroup.POST("/versions/:id/submit", templateHandler.SubmitForReview)
		templateGroup.POST("/versions/:id/approve", templateHandler.ApproveTemplate)
		templateGroup.POST("/versions/:id/reject", templateHandler.RejectTemplate)
		templateGroup.POST("/:template_id/versions/:version_id/rollback", templateHandler.RollbackTemplate)
		templateGroup.GET("/:template_id/content", templateHandler.GetTemplateContent)
	}

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
