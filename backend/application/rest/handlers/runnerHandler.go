package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/bersennaidoo/iathletestats/backend/domain/models"
	"github.com/bersennaidoo/iathletestats/backend/infrastructure/services"
	"github.com/gin-gonic/gin"
)

const ROLE_ADMIN = "admin"
const ROLE_RUNNER = "runner"

type RunnersHandler struct {
	runnersService *services.RunnersService
}

func NewRunnersHandler(runnersService *services.RunnersService) *RunnersHandler {
	return &RunnersHandler{
		runnersService: runnersService,
	}
}

func (rc RunnersHandler) CreateRunner(ctx *gin.Context) {

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
