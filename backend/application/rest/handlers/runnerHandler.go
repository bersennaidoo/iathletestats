package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/bersennaidoo/iathletestats/backend/domain/models"
	"github.com/bersennaidoo/iathletestats/backend/infrastructure/metrics"
	"github.com/bersennaidoo/iathletestats/backend/infrastructure/services"
	"github.com/gin-gonic/gin"
)

const ROLE_ADMIN = "admin"
const ROLE_RUNNER = "runner"

type RunnersHandler struct {
	runnersService *services.RunnersService
	usersService   *services.UsersService
}

func NewRunnersHandler(runnersService *services.RunnersService, usersService *services.UsersService) *RunnersHandler {
	return &RunnersHandler{
		runnersService: runnersService,
		usersService:   usersService,
	}
}

func (rc RunnersHandler) CreateRunner(ctx *gin.Context) {

	metrics.HttpRequestsCounter.Inc()

	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshaling create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rc.runnersService.CreateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
