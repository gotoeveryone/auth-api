package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gotoeveryone/general-api/app/config"
	"github.com/gotoeveryone/general-api/app/domain/entity"
	"github.com/gotoeveryone/general-api/app/infrastructure"
)

// Registration is execute account registration
func Registration(c *gin.Context) {
	// Execute validation
	var u entity.User
	if err := c.ShouldBindWith(&u, binding.JSON); err != nil {
		errorBadRequest(c, err.Error())
		return
	}

	// Check the same account already exists
	ur := infrastructure.NewUserRepository()
	if res, err := ur.Exists(u.Account); err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	} else if res {
		errorBadRequest(c, "Account is already exists")
		return
	}

	// Issue initial password
	password := config.Generate(16)

	// Create user
	if err := ur.Create(&u, password); err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"password": password,
	})
}

// Activate is enable account with update password
func Activate(c *gin.Context) {
	// Execute validation
	var a entity.Activate
	if err := c.ShouldBindWith(&a, binding.JSON); err != nil {
		errorBadRequest(c, err.Error())
		return
	}

	// Deny change to same password
	if a.Password == a.NewPassword {
		errorBadRequest(c, "Not allowed changing to same password")
		return
	}

	// Search user
	ur := infrastructure.NewUserRepository()
	user, err := ur.FindByUserAndPassword(a.Account, a.Password)
	if err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	}

	// Enable account with update password
	user.IsEnable = true
	if err := ur.UpdatePassword(user, a.NewPassword); err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	}

	c.JSON(http.StatusOK, user)
}

// Authenticate is execute user authenticate
func Authenticate(c *gin.Context) {
	// Execute validation
	var input entity.Authenticate
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		errorBadRequest(c, err.Error())
		return
	}

	// Search user
	ur := infrastructure.NewUserRepository()
	user, err := ur.FindByUserAndPassword(input.Account, input.Password)
	if err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	}

	// When initial password still not changed, Deny authentications
	if !user.IsActive {
		errorUnauthorized(c, "Password must be changed")
		return
	}

	// Deny disabled account
	if !user.IsEnable {
		errorUnauthorized(c, "Account is invalid")
		return
	}

	// Create token
	tr := infrastructure.NewTokenRepository()
	var token entity.Token
	if err := tr.Create(user, &token); err != nil {
		errorInternalServerError(c, err)
		return
	}

	// Authenticated
	if err := ur.UpdateAuthed(user); err != nil {
		errorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}

// GetUser is find user data from token
func GetUser(c *gin.Context) {
	// Find user from post token
	token := c.GetString(TokenKey)
	ur := infrastructure.NewUserRepository()
	user, err := ur.FindByToken(token)
	if err != nil {
		errorUnauthorized(c, "Authorization failed")
		return
	}

	c.JSON(http.StatusOK, user)
}

// Deauthenticate is execute user deauthentication
func Deauthenticate(c *gin.Context) {
	// Delete token
	token := c.GetString(TokenKey)
	tr := infrastructure.NewTokenRepository()
	if err := tr.Delete(token); err != nil {
		errorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
