package middleware

import jwt "github.com/appleboy/gin-jwt/v2"

// Auth is middleware interface for authentication and authorization
type Auth interface {
	Create() (*jwt.GinJWTMiddleware, error)
}
