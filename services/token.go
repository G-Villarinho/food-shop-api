package services

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/G-Villarinho/food-shop-api/config"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

//go:generate mockery --name=TokenService --output=../mocks --outpkg=mocks
type TokenService interface {
	CreateToken(userID uuid.UUID, sessionID uuid.UUID) (string, error)
	ExtractSessionID(token string) (uuid.UUID, error)
}

type tokenService struct {
	di *internal.Di
}

func NewTokenService(di *internal.Di) (TokenService, error) {
	return &tokenService{di: di}, nil
}

func (t *tokenService) CreateToken(userID uuid.UUID, sessionID uuid.UUID) (string, error) {
	privateKey, err := parseECPrivateKey(config.Env.PrivateKey)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"userId": userID.String(),
		"exp":    time.Now().Add(time.Duration(config.Env.Cache.SessionExp) * time.Hour).Unix(),
		"sid":    sessionID.String(),
		"iat":    time.Now().Unix(),
		"iss":    "level-up.com",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (t *tokenService) ExtractSessionID(token string) (uuid.UUID, error) {
	publicKey, err := parseECPublicKey(config.Env.PublicKey)
	if err != nil {
		return uuid.Nil, err
	}

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return publicKey, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	sid, ok := claims["sid"].(string)
	if !ok {
		return uuid.Nil, errors.New("session ID (sid) not found or invalid in token claims")
	}

	sessionID, err := uuid.Parse(sid)
	if err != nil {
		return uuid.Nil, err
	}
	return sessionID, nil
}

func parseECPrivateKey(pemKey string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("failed to parse EC private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func parseECPublicKey(pemKey string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to parse EC public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not a valid ECDSA public key")
	}

	return ecdsaPubKey, nil
}
