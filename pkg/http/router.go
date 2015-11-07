package http

import "github.com/gin-gonic/gin"

// RouterEngine is used to route HTTP traffic
type RouterEngine struct {
	*gin.Engine
}

// Router returns a configured router, with any additional middleware
func Router(m ...gin.HandlerFunc) RouterEngine {
	r := gin.New()

	for _, h := range m {
		r.Use(h)
	}

	return RouterEngine{r}
}
