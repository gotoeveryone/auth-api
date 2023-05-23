package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/domain/repository"
	"github.com/gotoeveryone/auth-api/app/presentation/handler"
	"github.com/rs/zerolog/log"
)

type userHandler struct {
	repo repository.User
}

// NewUserHandler is create action handler for user
func NewUserHandler(ur repository.User) handler.User {
	return &userHandler{
		repo: ur,
	}
}

// Register is execute registration of account
// @Summary Execute registration of account
// @Tags Authenticate
// @Accept  json
// @Produce json
// @Param data body entity.RegistrationUser true "request data"
// @Success 201 {object} entity.GeneratedPassword
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/users [post]
func (h *userHandler) Register(c *gin.Context) {
	var p entity.RegistrationUser
	if err := c.ShouldBindJSON(&p); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			errorBadRequest(c, ValidationErrors(verr, &p))
			return
		}
		errorBadRequest(c, errValidationFailed)
		return
	}

	// Check the same account already exists
	if res, err := h.repo.Exists(p.Account); err != nil {
		errorInternalServerError(c, err)
		return
	} else if res {
		errorBadRequest(c, errExistsAccount)
		return
	}

	t, err := time.Parse(p.Birthday, "2006-01-02")
	if err != nil {
		errorInternalServerError(c, err)
		return
	}

	u := entity.User{
		Account:     p.Account,
		Name:        p.Name,
		Gender:      entity.Gender(p.Gender),
		MailAddress: p.MailAddress,
		Birthday:    entity.Date{Time: t},
	}

	if p.Role != nil {
		u.Role = entity.Role(*p.Role)
	}

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
func (h *userHandler) Activate(c *gin.Context) {
	// Execute validation
	var a entity.Activate
	if err := c.ShouldBindJSON(&a); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			errorBadRequest(c, ValidationErrors(verr, &a))
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
		log.Fatal().Err(err)
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

// Identity is get authenticated user
// @Summary Return authenticated user
// @Tags Authenticate
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} entity.User
// @Failure 401 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 405 {object} entity.Error
// @Router /v1/me [get]
func (h *userHandler) Identity(c *gin.Context) {
	identity, _ := c.Get(config.IdentityKey)
	user := *identity.(*entity.User)

	c.JSON(http.StatusOK, user)
}
