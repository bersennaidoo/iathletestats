package services

import (
	"net/http"

	"github.com/bersennaidoo/iathletestats/backend/domain/models"
	"github.com/bersennaidoo/iathletestats/backend/foundation/hash"
	"github.com/bersennaidoo/iathletestats/backend/infrastructure/repositories/pgrepo"
)

type UsersService struct {
	usersRepository *pgrepo.UsersRepo
}

func NewUsersService(usersRepository *pgrepo.UsersRepo) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

func (us UsersService) Login(username string, password string) (string, *models.ResponseError) {
	if username == "" || password == "" {
		return "", &models.ResponseError{
			Message: "Invalid username or password",
			Status:  http.StatusBadRequest,
		}
	}

	id, responseErr := us.usersRepository.LoginUser(username, password)
	if responseErr != nil {
		return "", responseErr
	}

	if id == "" {
		return "", &models.ResponseError{
			Message: "Login failed",
			Status:  http.StatusUnauthorized,
		}
	}

	accessToken, responseErr := hash.GenerateAccessToken(username)
	if responseErr != nil {
		return "", responseErr
	}

	us.usersRepository.SetAccessToken(accessToken, id)

	return accessToken, nil
}

func (us UsersService) Logout(accessToken string) *models.ResponseError {
	if accessToken == "" {
		return &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}

	return us.usersRepository.RemoveAccessToken(accessToken)
}

func (us UsersService) AuthorizeUser(accessToken string, expectedRoles []string) (bool,
	*models.ResponseError) {
	if accessToken == "" {
		return false, &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}

	role, responseErr := us.usersRepository.GetUserRole(accessToken)
	if responseErr != nil {
		return false, responseErr
	}

	if role == "" {
		return false, &models.ResponseError{
			Message: "Failed to authorize user",
			Status:  http.StatusUnauthorized,
		}
	}

	for _, expectedRole := range expectedRoles {
		if expectedRole == role {
			return true, nil
		}
	}

	return false, nil
}
