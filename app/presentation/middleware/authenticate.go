package middleware

import "github.com/gin-gonic/gin"

// Auth is middleware interface for authentication and authorization
type Auth interface {
	Authorized() gin.HandlerFunc
}
