package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gotoeveryone/general-api/app/domain"
	"github.com/gotoeveryone/general-api/app/domain/entity"
	"github.com/gotoeveryone/general-api/app/infrastructure"
	"github.com/gotoeveryone/golib/logs"
	"gopkg.in/go-playground/validator.v8"
)

var (
	errInvalidAccount = errors.New("Account is invalid")
	errExistsAccount  = errors.New("Account is already exists")
	errInvalidRole    = errors.New("Role is invalid")

	errUpdatePassword     = errors.New("Update password failed")
	errSamePassword       = errors.New("Not allowed changing to same password")
	errMustChangePassword = errors.New("Password must be changed")
)

// Registration is execute account registration
func Registration(c *gin.Context) {
	// Execute validation
	var u entity.User
	if err := c.ShouldBindWith(&u, binding.JSON); err != nil {
		errors := err.(validator.ValidationErrors)
		ErrorBadRequest(c, domain.ValidationErrors(errors, &u))
		return
	}

	// Check the same account already exists
	ur := infrastructure.NewUserRepository()
	if res, err := ur.Exists(u.Account); err != nil {
		ErrorInternalServerError(c, err)
		return
	} else if res {
		ErrorBadRequest(c, errExistsAccount)
		return
	}

	// Check valid role
	if !ur.ValidRole(u.Role) {
		ErrorBadRequest(c, errInvalidRole)
		return
	}

	// Create user
	pass, err := ur.Create(&u)
	if err != nil {
		ErrorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"password": pass,
	})
}

// Activate is enable account with update password
func Activate(c *gin.Context) {
	// Execute validation
	var a entity.Activate
	if err := c.ShouldBindWith(&a, binding.JSON); err != nil {
		errors := err.(validator.ValidationErrors)
		ErrorBadRequest(c, domain.ValidationErrors(errors, &a))
		return
	}

	// Deny change to same password
	if a.Password == a.NewPassword {
		ErrorBadRequest(c, errSamePassword)
		return
	}

	// Search user
	ur := infrastructure.NewUserRepository()
	user, err := ur.FindByAccount(a.Account)
	if err != nil {
		ErrorInternalServerError(c, err)
		return
	}

	// Check password matching from user has password
	if err := ur.MatchPassword(user.Password, a.Password); err != nil {
		logs.Error(err)
		ErrorUnauthorized(c, ErrUnauthorized)
		return
	}

	// Enable account with update password
	user.IsEnable = true
	if err := ur.UpdatePassword(user, a.NewPassword); err != nil {
		ErrorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
	})
}

// Authenticate is execute user authenticate
func Authenticate(c *gin.Context) {
	// Execute validation
	var input entity.Authenticate
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		errors := err.(validator.ValidationErrors)
		ErrorBadRequest(c, domain.ValidationErrors(errors, &input))
		return
	}

	// Search user
	ur := infrastructure.NewUserRepository()
	user, err := ur.FindByAccount(input.Account)
	if err != nil {
		ErrorInternalServerError(c, err)
		return
	}

	// Invalid account
	if !ur.ValidUser(user) {
		ErrorBadRequest(c, errInvalidAccount)
		return
	}

	// When initial password still not changed, Deny authentications
	if !user.IsActive {
		ErrorBadRequest(c, errMustChangePassword)
		return
	}

	// Check password matching from user has password
	if err := ur.MatchPassword(user.Password, input.Password); err != nil {
		logs.Error(err)
		ErrorUnauthorized(c, ErrUnauthorized)
		return
	}

	// Create token
	tr := infrastructure.NewTokenRepository()
	var token entity.Token
	if err := tr.Create(user, &token); err != nil {
		ErrorInternalServerError(c, err)
		return
	}

	// Authenticated
	if err := ur.UpdateAuthed(user); err != nil {
		ErrorInternalServerError(c, err)
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
		ErrorInternalServerError(c, err)
		return
	}

	// Invalid account
	if !ur.ValidUser(user) {
		ErrorBadRequest(c, errInvalidAccount)
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
		ErrorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
