package handlers

import (
	"log"
	"net/http"

	"github.com/bersennaidoo/iathletestats/backend/infrastructure/services"
	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	usersService *services.UsersService
}

func NewUsersHandler(usersService *services.UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
	}
}

func (uc UsersHandler) Login(ctx *gin.Context) {
	username, password, ok := ctx.Request.BasicAuth()
	if !ok {
		log.Println("Error while reading credentials")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	accessToken, responseErr := uc.usersService.Login(username, password)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, accessToken)
}

func (uc UsersHandler) Logout(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")

	responseErr := uc.usersService.Logout(accessToken)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}
