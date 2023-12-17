package hash

import (
	"encoding/base64"
	"net/http"

	"github.com/bersennaidoo/iathletestats/backend/domain/models"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAccessToken(pwd string) (string, *models.ResponseError) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", &models.ResponseError{
			Message: "Failed to generate token",
			Status:  http.StatusInternalServerError,
		}
	}

	return base64.StdEncoding.EncodeToString(hash), nil
}
