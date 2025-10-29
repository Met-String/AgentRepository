package main

import (
	"log"
	// "net/http"
	"strings"
	"time"
	gateway_handler "github.com/Met-String/AgentSquare/internal/gateway/handler"
	extension_handler "github.com/Met-String/AgentSquare/internal/extension_observer/handler"
	resume_handler "github.com/Met-String/AgentSquare/internal/resume/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 网关 负责托管静态资源、Chat WS通讯。
func main() {
	r := gin.Default()

// 自定义配置
r.Use(cors.New(cors.Config{
	// AllowOrigins:     []string{"*"},  // 允许所有域
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowOriginFunc: func(origin string) bool {
		log.Println(origin)
		return strings.HasPrefix(origin, "chrome-extension://"+"nkkmhmcligmldeppkbcjfegoekfmclhl")
	},
	AllowCredentials: true,
	MaxAge: 12 * time.Hour,
}))

	// r.Use(func(c *gin.Context) {
	// 	user_agent := c.GetHeader("user-agent")
	// 	log.Println("user-agent:", user_agent)
	// })

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// WebSocket Echo 网关
	r.GET("/ws", gateway_handler.WsEchoHandler)
	// 前端静态资源
	gatewayAssets := gateway_handler.GatewayStaticHandler("./assets/gateway_vue3/dist", "/assets/gateway/")
	// r.GET("/", func(c *gin.Context) {
	// 	c.Redirect(http.StatusFound, "/assets/gateway/index.html")
	// })
	r.GET("/assets/gateway//*filepath", gatewayAssets)

	// ==========[简历]==========
	resumeAssets := resume_handler.ResumeStaticHandler("./assets/resume", "/resume/")
	r.GET("/resume/*filepath", resumeAssets)

	// ==========[浏览器拓展打点监听]==========
	r.POST("/extension", extension_handler.ExtensionEventHandler)

	addr := ":8080"
	log.Println("Gateway listening on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("server error:", err)
	}
}
