package controllers

import (
	"errors"
	"general-api/app/models"
	"general-api/app/services"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

const (
	// TokenPrefix ヘッダのトークンに付与されるPrefix
	TokenPrefix = "Bearer "
)

// State 状態
type State struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	LogLevel    string `json:"logLevel"`
	TimeZone    string `json:"timeZone"`
}

// GetState 状態
func GetState(c *gin.Context) {
	c.JSON(200, State{
		Status:      "OK",
		Environment: gin.Mode(),
		LogLevel:    services.Config.Log.Level,
		TimeZone:    time.Local.String(),
	})
}

// Authenticate 認証
func Authenticate(c *gin.Context) {
	// バリデーションエラー
	var postData models.Login
	// TODO: bindingでエラー発生時にHTMLが出力されてしまう
	if err := c.BindJSON(&postData); err != nil {
		ErrorJSON(c, 400, errors.New("Bad Request"))
		return
	}

	// ユーザの検索
	var userService services.UsersService
	user, err := userService.FindActiveUser(postData.Account)
	if err != nil {
		ErrorJSON(c, 401, errors.New("Credentials is invalid"))
		return
	}

	// パスワードの一致確認
	input := []byte(postData.Password)
	storeHashed := []byte(user.Password)
	if err := bcrypt.CompareHashAndPassword(storeHashed, input); err != nil {
		ErrorJSON(c, 401, errors.New("Credentials is invalid"))
		return
	}

	token := models.Token{
		UserID:      user.ID,
		Token:       generateToken(),
		Environment: gin.Mode(),
		Expire:      600,
	}
	var tokenService services.TokensService
	if err := tokenService.Create(token); err != nil {
		ErrorJSON(c, 500, err)
	}

	c.JSON(200, token)
}

// GetUser ユーザ取得
func GetUser(c *gin.Context) {
	// ヘッダに"Access-Token"が含まれているかチェック
	tokenHeader := c.Request.Header.Get("Access-Token")
	if !strings.HasPrefix(tokenHeader, TokenPrefix) {
		ErrorJSON(c, 401, errors.New("Token is required"))
		return
	}

	// ヘッダからトークンを取得
	token := strings.TrimLeft(tokenHeader, TokenPrefix)
	if token == "" {
		ErrorJSON(c, 401, errors.New("Token is required"))
		return
	}

	// トークンからユーザを取得
	var tokenService services.TokensService
	user, err := tokenService.FindUser(token)
	if err != nil {
		ErrorJSON(c, 401, errors.New("Token is invalid"))
		return
	}

	c.JSON(200, user)
}

// ErrorJSON エラー用JSONを出力します。
func ErrorJSON(c *gin.Context, code int, err error) {
	c.JSON(code, models.Error{
		Code:    code,
		Message: err.Error(),
	})
}

func generateToken() string {
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r := make([]byte, 50)
	for i := range r {
		r[i] = letters[rand.Intn(len(letters))]
	}
	return string(r)
}
