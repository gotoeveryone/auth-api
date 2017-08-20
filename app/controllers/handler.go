package controllers

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gotoeveryone/general-api/app/models"
	"github.com/gotoeveryone/general-api/app/services"
	"github.com/gotoeveryone/golib"
	"github.com/gotoeveryone/golib/logs"
)

const (
	// AuthorizedHeader 認可済みヘッダ
	AuthorizedHeader = "X-AUTHORIZED_USER"
)

// GetState 状態
func GetState(c *gin.Context) {
	c.JSON(200, models.State{
		Status:      "OK",
		Environment: gin.Mode(),
		LogLevel:    golib.AppConfig.Log.Level,
		TimeZone:    time.Local.String(),
	})
}

// Authenticate 認証
func Authenticate(c *gin.Context) {
	// バリデーション
	var postData models.Login
	if err := c.ShouldBindWith(&postData, binding.JSON); err != nil {
		errorJSON(c, 400, err)
		return
	}

	// ユーザの検索
	var userService services.UsersService
	user, err := userService.FindActiveUser(postData.Account)
	if err != nil {
		errorJSON(c, 401, err)
		return
	}

	// パスワードの一致確認
	input := []byte(postData.Password)
	storeHashed := []byte(user.Password)
	if err := bcrypt.CompareHashAndPassword(storeHashed, input); err != nil {
		errorJSON(c, 401, err)
		return
	}

	// トークンの生成
	var ts services.TokensService
	token := models.Token{
		UserID:      user.ID,
		Token:       ts.GenerateToken(),
		Environment: gin.Mode(),
		Expire:      600,
	}
	if err := ts.Create(token); err != nil {
		errorJSON(c, 500, err)
		return
	}

	c.JSON(200, token)
}

// GetUser ユーザ取得
func GetUser(c *gin.Context) {
	// トークンからユーザを取得
	token := c.Request.Header.Get(AuthorizedHeader)
	var ts services.TokensService
	user, err := ts.FindUser(token)
	if err != nil {
		errorJSON(c, 401, err)
		return
	}

	c.JSON(200, user)
}

// Deauthenticate 認証解除
func Deauthenticate(c *gin.Context) {
	// トークンの削除
	token := c.Request.Header.Get(AuthorizedHeader)
	var ts services.TokensService
	if err := ts.Delete(token); err != nil {
		errorJSON(c, 500, err)
		return
	}

	c.JSON(204, gin.H{})
}

// エラー用JSONの出力
func errorJSON(c *gin.Context, code int, err error) {
	logs.Error(err)
	var message string
	var header string
	switch code {
	case 400:
		message = "Bad Request"
		header = "error=\"invalid_request\""
	case 401:
		message = "Authenticate Failed"
		header = "error=\"invalid_token\""
	case 500:
	default:
		message = "Internal Server Error"
	}
	if header != "" && c.Request.Header.Get(AuthorizedHeader) != "" {
		c.Writer.Header().Set("WWW-Authenticate", "Bearer "+header)
	}
	c.AbortWithStatusJSON(code, models.Error{
		Code:    code,
		Message: message,
	})
}
