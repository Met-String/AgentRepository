package handler
import (
	"net/http"
	"github.com/gin-gonic/gin"
)


// 返回一个 Gin 处理器，用于将指定目录作为静态文件服务器暴露。
func ResumeStaticHandler(root string, prefix string) gin.HandlerFunc {
	fileServer := http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
	return func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}