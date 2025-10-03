package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/Met-String/AgentSquare/internal/gateway/handler"
)

func main() {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// WebSocket Echo 网关
	r.GET("/ws", handler.WsEchoHandler)

	addr := ":8080"
	log.Println("Gateway listening on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server error:", err)
	}
}
