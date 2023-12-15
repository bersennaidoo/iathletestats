// Package models represents the domain models for the iathletestats program.
package models

// ResponseError is used to track and record io errors conditions for the program.
type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"-"`
}
