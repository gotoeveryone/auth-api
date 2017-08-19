package middlewares

import (
	"general-api/app/controllers"
	"general-api/app/models"
	"general-api/app/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotoeveryone/golib/logs"
)

const (
	// ヘッダのトークンに付与されるPrefix
	tokenPrefix = "Bearer "
)

// HasToken トークン保持確認
func HasToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// "Authorization"ヘッダを含むか
		tokenHeader := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(tokenHeader, tokenPrefix) {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer realm=\"token_required\"")
			c.AbortWithStatusJSON(401, models.Error{
				Code:    401,
				Message: "Token is required",
			})
			return
		}

		// ヘッダからトークンが取得できるか
		token := strings.TrimSpace(strings.Replace(tokenHeader, tokenPrefix, "", 1))
		if token == "" {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer error=\"token_required\"")
			c.AbortWithStatusJSON(401, models.Error{
				Code:    401,
				Message: "Token is required",
			})
			return
		}

		// トークンが不正ではないか
		var ts services.TokensService
		var m models.Token
		if err := ts.FindToken(token, &m); err != nil {
			logs.Error(err)
			c.Writer.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
			c.AbortWithStatusJSON(401, models.Error{
				Code:    401,
				Message: "Token is invalid",
			})
			return
		}

		// 認可済みヘッダを設定
		c.Request.Header.Set(controllers.AuthorizedHeader, m.Token)
		c.Next()
	}
}
