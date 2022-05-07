package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/presentation/handler"
	"github.com/sirupsen/logrus"
)

type authHandler struct {
	repo repository.User
}

// NewAuthHandler is create action handler for auth
func NewAuthHandler(ur repository.User) handler.Authenticate {
	return &authHandler{
		repo: ur,
	}
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
	if res, err := h.repo.Exists(u.Account); err != nil {
		errorInternalServerError(c, err)
		return
	} else if res {
		errorBadRequest(c, errExistsAccount)
		return
	}

	// Check valid role
	if !u.ValidRole() {
		errorBadRequest(c, errInvalidRole)
		return
	}

	// Create user
	pass, err := h.repo.Create(&u)
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
	user, err := h.repo.FindByAccount(a.Account)
	if err != nil {
		errorInternalServerError(c, err)
		return
	}
	if user == nil {
		errorUnauthorized(c, errUnauthorized)
		return
	}

	// Check password matching from user has password
	if err := h.repo.MatchPassword(user.Password, a.Password); err != nil {
		logrus.Error(err)
		errorUnauthorized(c, errUnauthorized)
		return
	}

	// Enable account with update password
	user.IsEnable = true
	if err := h.repo.UpdatePassword(user, a.NewPassword); err != nil {
		errorInternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// User is find user data from token
// @Summary Return authenticated user
// @Tags Authenticate
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} entity.User
// @Failure 401 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/me [get]
func (h *authHandler) User(c *gin.Context) {
	identity, _ := c.Get(config.IdentityKey)
	user := *identity.(*entity.User)

	c.JSON(http.StatusOK, user)
}
