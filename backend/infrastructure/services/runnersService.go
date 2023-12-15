package services

import (
	"net/http"

	"github.com/bersennaidoo/iathletestats/backend/domain/models"
	"github.com/bersennaidoo/iathletestats/backend/infrastructure/repositories/pgrepo"
)

type RunnersService struct {
	runnersRep *pgrepo.RunnersRepo
}

func NewRunnersService(runnersRep *pgrepo.RunnersRepo) *RunnersService {
	return &RunnersService{
		runnersRep: runnersRep,
	}
}

func (rs RunnersService) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunner(runner)
	if responseErr != nil {
		return nil, responseErr
	}

	return rs.runnersRepository.CreateRunner(runner)
}

func validateRunner(runner *models.Runner) *models.ResponseError {
	if runner.FirstName == "" {
		return &models.ResponseError{
			Message: "Invalid first name",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.LastName == "" {
		return &models.ResponseError{
			Message: "Invalid last name",
			Status:  http.StatusBadRequest}
	}

	if runner.Age < 0 || runner.Age > 125 {
		return &models.ResponseError{
			Message: "Invalid age",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.Country == "" {
		return &models.ResponseError{
			Message: "Invalid country",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}

func validateRunnerId(runnerId string) *models.ResponseError {
	if runnerId == "" {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}
