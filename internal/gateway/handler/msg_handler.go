package handler

import (
	"log"
	"net/http"
	"sync"
	"time"
	"github.com/Met-String/AgentSquare/internal/gateway/hub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


// WebSocket 升级器（演示用：放开跨域；生产环境请按需校验 Origin）
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handler.go (只保留对 Hub 的调用)
// 在包级别持有一个 hub
var ws_hub = hub.NewHub()

// WsEchoHandler 将 HTTP 连接升级为 WebSocket，并实现“打印并原样回显”
func WsEchoHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// 升级失败通常是 400/404/握手错误
		log.Println("upgrade error:", err)
		return
	}

	ws_hub.AddConn(conn)
	defer func() {
		_ = conn.Close()
		ws_hub.RemoveConn(conn)
	}()
	log.Printf("ws connected: %s", conn.RemoteAddr())

	// 基本的超时与心跳参数配置
	const (
		readLimit = 1 << 20 // 1MB
		pongWait  = 60 * time.Second
		writeWait = 10 * time.Second
		heartBeat = 30 * time.Second
	)
	conn.SetReadLimit(readLimit)

	// 持续Ping、Pong以维持连接
	// 接收Pong
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { 
		_ = conn.SetReadDeadline(time.Now().Add(pongWait)) // 收到客户端 Pong 就延长读取期限
		return nil
	})
	// 发出Ping
	wrLock := sync.Mutex{}
	pingTicker := time.NewTicker(heartBeat)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-pingTicker.C:
				wrLock.Lock()
				if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
					log.Println("ping error:", err)
					wrLock.Unlock()
					return
				}
				wrLock.Unlock()
			case <-done:
				return
			}
		}
	}()

	// 简单读写回环：读到什么就打印并原样写回
	for {
		msgType, msg, err := conn.ReadMessage() // 阻塞读
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

		// 在读循环中：
		if msgType == websocket.TextMessage {
			ws_hub.Broadcast(msg)
		}
	}

	close(done)
	pingTicker.Stop()
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