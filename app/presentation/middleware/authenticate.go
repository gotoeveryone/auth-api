package middleware

import "github.com/gin-gonic/gin"

// Authenticate is authenticate middleware
type Authenticate interface {
	Authorized() gin.HandlerFunc
}
