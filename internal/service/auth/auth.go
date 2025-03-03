package auth

import (
	"github.com/gin-gonic/gin"
)

func GignInWithProvider(c *gin.Context) {
	provider := c.Param("provider")

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	
}