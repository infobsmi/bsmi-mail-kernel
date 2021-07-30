package http_utils

import "github.com/gin-gonic/gin"

func GetClientIp(c *gin.Context) string {
	clientIp := c.GetHeader("HTTP_CF_CONNECTING_IP")
	if clientIp != "" {
		return clientIp
	}
	clientIp = c.GetHeader("X-Forwarded-For")
	if clientIp != "" {
		return clientIp
	}
	return c.Request.RemoteAddr
}
