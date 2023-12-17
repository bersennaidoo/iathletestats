package pgrepo

import (
	"database/sql"
	"net/http"

	"github.com/bersennaidoo/iathletestats/backend/domain/models"
)

type UsersRepo struct {
	dbClient *sql.DB
}

func NewUsersRepo(dbClient *sql.DB) *UsersRepo {
	return &UsersRepo{
		dbClient: dbClient,
	}
}

func (ur UsersRepo) LoginUser(username string, password string) (string,
	*models.ResponseError) {

	query := `
          SELECT id
		FROM users
		WHERE username = $1 and user_password = crypt($2, user_password)`

	rows, err := ur.dbClient.Query(query, username, password)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var id string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return id, nil
}

func (ur UsersRepo) GetUserRole(accessToken string) (string, *models.ResponseError) {
	query := `
		SELECT user_role
		FROM users
		WHERE access_token = $1`

	rows, err := ur.dbClient.Query(query, accessToken)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var role string
	for rows.Next() {
		err := rows.Scan(&role)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return role, nil
}

func (ur UsersRepo) SetAccessToken(accessToken string, id string) *models.ResponseError {
	query := `UPDATE users SET access_token = $1 WHERE id = $2`

	_, err := ur.dbClient.Exec(query, accessToken, id)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (ur UsersRepo) RemoveAccessToken(accessToken string) *models.ResponseError {
	query := `UPDATE users SET access_token = '' WHERE access_token = $1`

	_, err := ur.dbClient.Exec(query, accessToken)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}
