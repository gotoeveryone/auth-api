package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gotoeveryone/general-api/app/models"
	"github.com/gotoeveryone/general-api/app/services"
	"github.com/gotoeveryone/general-api/app/utils"
)

// Publish ユーザ登録
func Publish(c *gin.Context) {
	// バリデーション
	var u models.User
	if err := c.ShouldBindWith(&u, binding.JSON); err != nil {
		errorBadRequest(c, err.Error())
		return
	}

	// 同じアカウントのユーザがすでに存在するか
	var us services.UsersService
	if res, err := us.Exists(u.Account); err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	} else if res {
		errorBadRequest(c, "Account is already exists")
		return
	}

	// 初期パスワードの発行
	password := utils.Generate(16)

	// 一般ユーザとして登録
	if err := us.Create(&u, password); err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"password": password,
	})
}

// Activate アカウント有効化
func Activate(c *gin.Context) {
	// バリデーション
	var a models.Activate
	if err := c.ShouldBindWith(&a, binding.JSON); err != nil {
		errorBadRequest(c, err.Error())
		return
	}

	// 同じパスワードには変更させない
	if a.Password == a.NewPassword {
		errorBadRequest(c, "Not allowed changing to same password")
		return
	}

	// ユーザの検索
	var us services.UsersService
	user, err := us.FindUser(a.Account, a.Password)
	if err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	}

	// アカウントを有効化し、パスワードを更新
	user.IsEnable = true
	if err := us.UpdatePassword(user, a.NewPassword); err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	}

	c.JSON(http.StatusOK, user)
}

// Authenticate 認証
func Authenticate(c *gin.Context) {
	// バリデーション
	var input models.Login
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		errorBadRequest(c, err.Error())
		return
	}

	// ユーザの検索
	var us services.UsersService
	user, err := us.FindUser(input.Account, input.Password)
	if err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	}

	// パスワード変更未実施
	if !user.IsActive {
		errorUnauthorized(c, "Password must be changed")
		return
	}

	// 無効アカウント
	if !user.IsEnable {
		errorUnauthorized(c, "Account is invalid")
		return
	}

	// トークンの生成
	var ts services.TokensService
	var token models.Token
	if err := ts.Create(user, &token); err != nil {
		errorInternalServerError(c, err)
		return
	}

	// 認証日時を更新
	if err := us.UpdateAuthed(user); err != nil {
		errorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}

// GetUser ユーザ取得
func GetUser(c *gin.Context) {
	// トークンからユーザを取得
	token := c.GetString(TokenKey)
	var ts services.TokensService
	user, err := ts.FindUser(token)
	if err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	}

	c.JSON(http.StatusOK, user)
}

// Deauthenticate 認証解除
func Deauthenticate(c *gin.Context) {
	// トークンの削除
	token := c.GetString(TokenKey)
	var ts services.TokensService
	if err := ts.Delete(token); err != nil {
		errorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
