package utils

import (
	"github.com/mjaliz/deviran/internal/constants"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), constants.HashPasswordCost)
	return string(bytes), err
}

func CompareHashAndPass(userPass, payloadPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userPass), []byte(payloadPass))
	if err != nil {
		return err
	}
	return nil
}
