package services

import (
	"time"

	"github.com/G-Villarinho/level-up-api/config"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	CreateToken(userID uuid.UUID, sessionID uuid.UUID, restaurantID *uuid.UUID) (string, error)
	ExtractSessionID(token string) (uuid.UUID, error)
}

type tokenService struct {
	di *internal.Di
}

func NewTokenService(di *internal.Di) (TokenService, error) {
	return &tokenService{di: di}, nil
}

func (t *tokenService) CreateToken(userID uuid.UUID, sessionID uuid.UUID, restaurantID *uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID.String(),
		"exp":    time.Now().Add(time.Duration(config.Env.Cache.SessionExp) * time.Hour).Unix(),
		"sid":    sessionID.String(),
		"iat":    time.Now().Unix(),
		"iss":    "level-up.com",
	}

	if restaurantID != nil {
		claims["rid"] = restaurantID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedToken, err := token.SignedString(config.Env.PrivateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (t *tokenService) ExtractSessionID(token string) (uuid.UUID, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return config.Env.PublicKey, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	sessionID, err := uuid.Parse(claims["sid"].(string))
	if err != nil {
		return uuid.Nil, err
	}

	return sessionID, nil
}
