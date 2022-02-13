package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gotoeveryone/auth-api/app/domain"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/presentation/handler"
	"github.com/sirupsen/logrus"
)

type (
	// authenticateHandler is authentication action handler
	authenticateHandler struct {
		userRepo  repository.UserRepository
		tokenRepo repository.TokenRepository
	}
)

// NewAuthenticateHandler is state action handler
func NewAuthenticateHandler(ur repository.UserRepository, tr repository.TokenRepository) handler.AuthenticateHandler {
	return &authenticateHandler{
		userRepo:  ur,
		tokenRepo: tr,
	}
}

// Authenticate is execute user authenticate
func (h *authenticateHandler) Authenticate(c *gin.Context) {
	// Execute validation
	var input entity.Authenticate
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			errorBadRequest(c, domain.ValidationErrors(verr, &input))
			return
		}
		errorBadRequest(c, errValidationFailed)
		return
	}

	// Search user
	user, err := h.userRepo.FindByAccount(input.Account)
	if err != nil {
		errorInternalServerError(c, err)
		return
	}

	// Invalid account
	if !h.userRepo.ValidUser(user) {
		errorBadRequest(c, errInvalidAccount)
		return
	}

	// When initial password still not changed, Deny authentications
	if !user.IsActive {
		errorBadRequest(c, errMustChangePassword)
		return
	}

	// Check password matching from user has password
	if err := h.userRepo.MatchPassword(user.Password, input.Password); err != nil {
		logrus.Error(err)
		errorUnauthorized(c, errUnauthorized)
		return
	}

	// Create token
	var token entity.Token
	if err := h.tokenRepo.Create(user, &token); err != nil {
		errorInternalServerError(c, err)
		return
	}

	// Authenticated
	if err := h.userRepo.UpdateAuthed(user); err != nil {
		errorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}

// Deauthenticate is execute user deauthentication
func (h *authenticateHandler) Deauthenticate(c *gin.Context) {
	// Delete token
	token := c.GetString(TokenKey)
	if err := h.tokenRepo.Delete(token); err != nil {
		errorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// Registration is execute account registration
func (h *authenticateHandler) Registration(c *gin.Context) {
	// Execute validation
	var u entity.User
	if err := c.ShouldBindWith(&u, binding.JSON); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			errorBadRequest(c, domain.ValidationErrors(verr, &u))
			return
		}
		errorBadRequest(c, errValidationFailed)
		return
	}

	// Check the same account already exists
	if res, err := h.userRepo.Exists(u.Account); err != nil {
		errorInternalServerError(c, err)
		return
	} else if res {
		errorBadRequest(c, errExistsAccount)
		return
	}

	// Check valid role
	if !h.userRepo.ValidRole(u.Role) {
		errorBadRequest(c, errInvalidRole)
		return
	}

	// Create user
	pass, err := h.userRepo.Create(&u)
	if err != nil {
		errorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"password": pass,
	})
}

// Activate is enable account with update password
func (h *authenticateHandler) Activate(c *gin.Context) {
	// Execute validation
	var a entity.Activate
	if err := c.ShouldBindWith(&a, binding.JSON); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			errorBadRequest(c, domain.ValidationErrors(verr, &a))
			return
		}
		errorBadRequest(c, errValidationFailed)
		return
	}

	// Deny change to same password
	if a.Password == a.NewPassword {
		errorBadRequest(c, errSamePassword)
		return
	}

	// Search user
	user, err := h.userRepo.FindByAccount(a.Account)
	if err != nil {
		errorInternalServerError(c, err)
		return
	}

	// Check password matching from user has password
	if err := h.userRepo.MatchPassword(user.Password, a.Password); err != nil {
		logrus.Error(err)
		errorUnauthorized(c, errUnauthorized)
		return
	}

	// Enable account with update password
	user.IsEnable = true
	if err := h.userRepo.UpdatePassword(user, a.NewPassword); err != nil {
		errorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
	})
}

// GetUser is find user data from token
func (h *authenticateHandler) GetUser(c *gin.Context) {
	// Find user from post token
	token := c.GetString(TokenKey)
	var t entity.Token
	if err := h.tokenRepo.Find(token, &t); err != nil {
		errorInternalServerError(c, err)
		return
	}

	var u entity.User
	if err := h.userRepo.Find(t.UserID, &u); err != nil {
		errorInternalServerError(c, err)
		return
	}

	// Invalid account
	if !h.userRepo.ValidUser(&u) {
		errorBadRequest(c, errInvalidAccount)
		return
	}

	c.JSON(http.StatusOK, u)
}
