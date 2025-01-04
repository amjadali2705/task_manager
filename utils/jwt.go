package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

const secretKey = "supersecret"
const secretKey2 = "supertopsecret"

func GenerateJwtToken(userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Minute * 1).Unix(),
	})

	Logger.Info("User Token generated successfully")
	return token.SignedString([]byte(secretKey))
}

func GenerateRefreshToken(userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	Logger.Info("Refresh Token generated successfully")
	return token.SignedString([]byte(secretKey2))
}

func VerifyJwtToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			Logger.Error("Unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		Logger.Error("Failed to parse token", zap.Error(err))
		return 0, errors.New("could not parse the token")
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		Logger.Error("Invalid token")
		return 0, errors.New("invalid Token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		Logger.Error("Invalid token claims")
		return 0, errors.New("invalid token claims")
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		Logger.Error("Invalid user id in token claims")
		return 0, errors.New("invalid token claims")
	}

	Logger.Info("User Token verified successfully", zap.Int64("userId", int64(userId)))
	return int64(userId), nil
}

func VerifyRefreshToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			Logger.Error("Unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secretKey2), nil
	})

	if err != nil {
		Logger.Error("Failed to parse token", zap.Error(err))
		return 0, errors.New("could not parse the token")
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		Logger.Error("Invalid token")
		return 0, errors.New("invalid Token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		Logger.Error("Invalid token claims")
		return 0, errors.New("invalid token claims")
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		Logger.Error("Invalid user id in token claims")
		return 0, errors.New("invalid token claims")
	}

	Logger.Info("Refresh Token verified successfully", zap.Int64("userId", int64(userId)))
	return int64(userId), nil
}

func DecodeJwtToken(token string) (map[string]interface{}, error) {
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		Logger.Error("Error decoding token", zap.Error(err))
		return nil, errors.New("could not decode the token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		Logger.Error("Invalid token claims")
		return nil, errors.New("invalid token claims")
	}

	claimsMap := make(map[string]interface{})
	for key, value := range claims {
		claimsMap[key] = value
	}

	Logger.Info("Token decoded successfully")
	return claimsMap, nil
}
