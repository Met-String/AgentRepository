package hub

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type Hub struct { 
	conns     map[*websocket.Conn]bool
	connsMu   sync.RWMutex
	broadcast chan []byte // 这他妈什么？一个channel？传递的是二进制字节流数据？
}

// NewHub initialises the hub and starts the broadcast loop.
func NewHub() *Hub { // 新建一个中心，而且返回值又是一个指针？我怀疑在GoLang里面，指针被习惯当作对象用。
	h := &Hub{ // 新建了一个Hub对象，并且当场取其地址？太带派了。
		conns:     make(map[*websocket.Conn]bool), // 我仍然不是那么习惯Go的结构体初始化方式...
		broadcast: make(chan []byte, 256), // 很好，我们有了256字节的缓存长度
	}
	go h.run() // go类型最糟糕的一点就是，类型的定义和方法的添加常常完全在不同的地方
	return h
}

func (h *Hub) run() { 
	for msg := range h.broadcast { // 好，run，然后呢？h.broadcast里什么都没有！
		h.connsMu.RLock() // Read锁也是上上了
		for c := range h.conns {// 遍历所有连接，嗯...接下来像是要广播的样子。
			go func(conn *websocket.Conn, m []byte) { // 我第一次见到这种协程的写法...
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) //🟢5 //没用的WriteDeadline
				if err := conn.WriteMessage(websocket.TextMessage, m); err != nil {
					log.Println("broadcast error:", err)
				}
			}(c, msg)
		}
		h.connsMu.RUnlock()
	}
}

func (h *Hub) AddConn(c *websocket.Conn) {
	h.connsMu.Lock()
	h.conns[c] = true
	h.connsMu.Unlock()
}

// RemoveConn drops the socket from the hub and closes it.
func (h *Hub) RemoveConn(c *websocket.Conn) {
	h.connsMu.Lock()
	delete(h.conns, c)
	h.connsMu.Unlock()
	_ = c.Close()
}

// Broadcast schedules the message to be sent to all active connections.
func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}
