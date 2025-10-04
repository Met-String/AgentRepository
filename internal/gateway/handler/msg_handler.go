package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// WebSocket 升级器（演示用：放开跨域；生产环境请按需校验 Origin）
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsEchoHandler 将 HTTP 连接升级为 WebSocket，并实现“打印并原样回显”
func WsEchoHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// 升级失败通常是 400/404/握手错误
		log.Println("upgrade error:", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	// 基本的超时与心跳设置（可选，稳定连接）
	const (
		readLimit = 1 << 20 // 1MB
		pongWait  = 60 * time.Second
		writeWait = 10 * time.Second
	)
	conn.SetReadLimit(readLimit)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		// 收到客户端 Pong 就延长读取期限
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	log.Printf("ws connected: %s", conn.RemoteAddr())

	// 简单读写回环：读到什么就打印并原样写回
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			// 客户端关闭/网络错误
			log.Println("read error:", err)
			break
		}

		// 打印收到的消息
		if msgType == websocket.TextMessage {
			log.Printf("recv (text) from %s: %q", conn.RemoteAddr(), string(msg))
		} else {
			log.Printf("recv (binary) from %s: %d bytes", conn.RemoteAddr(), len(msg))
		}

		// 回显
		_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println("write error:", err)
			break
		}
	}
	log.Printf("ws closed: %s", conn.RemoteAddr())
}

// GatewayStaticHandler 返回一个 Gin 处理器，用于将指定目录作为静态文件服务器暴露
func GatewayStaticHandler(root string) gin.HandlerFunc {
	fileServer := http.FileServer(http.Dir(root))
	return func(c *gin.Context) {
		requestedPath := c.Param("filepath")
		originalPath := c.Request.URL.Path
		c.Request.URL.Path = requestedPath
		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Request.URL.Path = originalPath
	}
}
