package main

import (
	"log"
	"net/http"
	"github.com/Met-String/AgentSquare/internal/gateway/handler"
	"github.com/gin-gonic/gin"
)

// 网关 负责托管静态资源、Chat WS通讯。
func main() {
	r := gin.Default()
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	// WebSocket Echo 网关
	r.GET("/ws", handler.WsEchoHandler)
	// 前端静态资源
	gatewayAssets := handler.GatewayStaticHandler("./assets/gateway")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/assets/gateway/index.html")
	})
	r.GET("/assets/gateway/*filepath", gatewayAssets)
	addr := ":8080"
	log.Println("Gateway listening on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server error:", err)
	}
}
