package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func DeviceTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		deviceType := getDeviceType(userAgent)

		// 将设备类型存储到 Gin 的上下文中
		c.Set("device_type", deviceType)

		// 继续处理请求
		c.Next()
	}
}

func getDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		return "phone"
	}

	return "pc"
}
