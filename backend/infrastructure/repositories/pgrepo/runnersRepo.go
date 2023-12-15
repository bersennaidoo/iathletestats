// Package pgrepo implements storage for storing domain models in postgresql.
package pgrepo

import (
	"database/sql"
	"net/http"

	"github.com/bersennaidoo/iathletestats/backend/domain/models"
)

// RunnersRepo represents a type to perform crud operations for runners.
type RunnersRepo struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

// NewRunnersRepo creates an instance of RunnersRepo initialized with a dbHandler.
func NewRunnersRepo(dbHandler *sql.DB) *RunnersRepo {
	return &RunnersRepo{
		dbHandler: dbHandler,
	}
}

// CreateRunner creates and returns a runner.
func (rr RunnersRepo) CreateRunner(runner *models.Runner) (*models.Runner,
	*models.ResponseError) {

	query := `
		INSERT INTO runners(first_name, last_name, age, country)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	// Query runs query and return table rows.
	rows, err := rr.dbHandler.Query(query, runner.FirstName, runner.LastName,
		runner.Age, runner.Country)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var runnerId string

	// rows.Next check for rows and if true then rows.Scan fetches id and parses into runnerId.
	// The last runnerId scanned identifies the runner insert into row.
	for rows.Next() {
		err := rows.Scan(&runnerId)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	// Constructs runner and returns to caller.
	return &models.Runner{
		ID:        runnerId,
		FirstName: runner.FirstName,
		LastName:  runner.LastName,
		Age:       runner.Age,
		IsActive:  true,
		Country:   runner.Country,
	}, nil
}
