package router

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginRouter struct {
	dispatcher *gin.Engine
}

func NewGinRouter() Router {
	return &ginRouter{dispatcher: gin.Default()}
}

func (r *ginRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.dispatcher.GET(uri, func(c *gin.Context) {
		ginWrapper(f, c)
	})
}

func (r *ginRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.dispatcher.POST(uri, func(c *gin.Context) {
		ginWrapper(f, c)
	})
}

func (r *ginRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.dispatcher.PUT(uri, func(c *gin.Context) {
		ginWrapper(f, c)
	})
}

func (r *ginRouter) SERVE(port string) error {
	log.Printf("Gin HTTP server running on port %v", port)
	return r.dispatcher.Run(":" + port)
}

// ginWrapper adapts the standard HTTP handler function to work with Gin's context.
func ginWrapper(f func(w http.ResponseWriter, r *http.Request), c *gin.Context) {
	// Convert Gin context to standard http.Request
	r := c.Request
	w := c.Writer

	// Extract params and inject into context
	params := ginExtractParams(c)
	ctx := r.Context()
	for k, v := range params {
		ctx = context.WithValue(ctx, k, v)
	}

	// Call the wrapped handler with modified context
	f(w, r.WithContext(ctx))
}

// ginExtractParams extracts URL parameters in a format compatible with `router.go`.
func ginExtractParams(c *gin.Context) map[string]string {
	return map[string]string{
		"gameName": c.Param("gamename"),
		"userName": c.Param("username"),
	}
}
