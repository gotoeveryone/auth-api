package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gotoeveryone/auth-api/app/domain"
	"github.com/gotoeveryone/auth-api/app/domain/entity"
	"github.com/gotoeveryone/auth-api/app/infrastructure"
	validator "gopkg.in/go-playground/validator.v8"
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
