package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS cross-domain middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //Request header
		if origin != "" {
			// Receive the origin sent by the client (important!)
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// All cross-domain request methods supported by the server
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			// Allow cross-domain settings to return other sub-segments and customize fields
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// Allow headers that the browser (client) can parse (important)
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			// Set cache time
			c.Header("Access-Control-Max-Age", "172800")
			// Allow the client to pass verification information such as cookies (important)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Allow type checking
		if method == http.MethodOptions {
			c.AbortWithStatus(200)
			return
		}

		// defer func() {
		// 	if err := recover(); err != nil {
		// 		log.Printf("Panic info is: %v", err)
		// 	}
		// }()

		c.Next()
	}
}
