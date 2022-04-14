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

type authHandler struct {
	userRepo  repository.User
	tokenRepo repository.Token
}

// NewAuthHandler is create action handler for auth
func NewAuthHandler(ur repository.User, tr repository.Token) handler.Authenticate {
	return &authHandler{
		userRepo:  ur,
		tokenRepo: tr,
	}
}

// Authenticate is execute authentication by user
// @Summary Execute authentication by user
// @Tags Authenticate
// @Accept  json
// @Produce json
// @Param data body entity.Authenticate true "request data"
// @Success 200 {object} entity.Token
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/auth [post]
func (h *authHandler) Authenticate(c *gin.Context) {
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

// Deauthenticate is execute deauthentication by user
// @Summary Execute deauthentication by user
// @Tags Authenticate
// @Security ApiKeyAuth
// @Produce json
// @Success 204
// @Failure 400 {object} entity.Error
// @Failure 401 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/deauth [delete]
func (h *authHandler) Deauthenticate(c *gin.Context) {
	// Delete token
	token := c.GetString(TokenKey)
	if err := h.tokenRepo.Delete(token); err != nil {
		errorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// Registration is execute registration of account
// @Summary Execute registration of account
// @Tags Authenticate
// @Accept  json
// @Produce json
// @Param data body entity.User true "request data"
// @Success 201 {object} entity.GeneratedPassword
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/users [post]
func (h *authHandler) Registration(c *gin.Context) {
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

	c.JSON(http.StatusCreated, entity.GeneratedPassword{
		Password: pass,
	})
}

// Activate is enable account with update password
// @Summary Enable account with update password
// @Tags Authenticate
// @Accept  json
// @Produce json
// @Param data body entity.Activate true "request data"
// @Success 200
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/activate [post]
func (h *authHandler) Activate(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{})
}

// GetUser is find user data from token
// @Summary Return authenticated user
// @Tags Authenticate
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} entity.User
// @Failure 401 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/users [get]
func (h *authHandler) GetUser(c *gin.Context) {
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
