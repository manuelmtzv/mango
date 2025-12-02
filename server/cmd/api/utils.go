package main

import (
	"github.com/manuelmtzv/mangocatnotes-api/internal/auth"
)

func HashPassword(password string) (string, error) {
	return auth.HashPassword(password)
}

func VerifyPassword(password, hash string) (bool, error) {
	return auth.VerifyPassword(password, hash)
}
