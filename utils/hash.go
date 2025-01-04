package utils

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		Logger.Error("Failed to hash the password", zap.Error(err))
		return "", err
	}

	Logger.Info("Password hashed successfully")
	return string(bytes), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		Logger.Error("Password hash check failed", zap.Error(err))
		return false
	}

	Logger.Info("Password hash check successful")
	return true
}
