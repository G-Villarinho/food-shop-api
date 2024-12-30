package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/G-Villarinho/food-shop-api/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTokenService_CreateToken(t *testing.T) {
	config.Env.PrivateKey = `-----BEGIN PRIVATE KEY-----
								MIGEAgEAMBAGByqGSM49AgEGBSuBBAAKBG0wawIBAQQg5RehdfX0dB5a1b4O9LRW
								GrasAxpShIkcufp95TfCCzyhRANCAASvtRoBafzVhFxPK2hnlLDrxIagEr/Eea7/
								37M6Iqt9tfeF03AJcIAVw/QmtI0yuxsKoWMYbnlJoNtCwrvJ8HD0
								-----END PRIVATE KEY-----`
	config.Env.PublicKey = `-----BEGIN PUBLIC KEY-----
								MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEr7UaAWn81YRcTytoZ5Sw68SGoBK/xHmu
								/9+zOiKrfbX3hdNwCXCAFcP0JrSNMrsbCqFjGG55SaDbQsK7yfBw9A==
								-----END PUBLIC KEY-----`
	config.Env.Cache.SessionExp = 1

	tokenService := &tokenService{}

	t.Run("should create a valid token successfully", func(t *testing.T) {
		userID := uuid.New()
		sessionID := uuid.New()

		token, err := tokenService.CreateToken(userID, sessionID)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
			return config.Env.PublicKey, nil
		})

		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)
		assert.Equal(t, userID.String(), claims["userId"])
		assert.Equal(t, sessionID.String(), claims["sid"])
		assert.Equal(t, "level-up.com", claims["iss"])
		assert.NotNil(t, claims["iat"])
		assert.NotNil(t, claims["exp"])
	})

	t.Run("should return an error when signing fails", func(t *testing.T) {
		// Simula chave inválida ao redefinir a assinatura do método `SignedString`
		originalPrivateKey := config.Env.PrivateKey
		defer func() {
			config.Env.PrivateKey = originalPrivateKey // Restaura a chave válida após o teste
		}()

		// Define uma chave privada inválida como substituta temporária
		config.Env.PrivateKey = "private_key" // Estrutura não inicializada

		userID := uuid.New()
		sessionID := uuid.New()

		token, err := tokenService.CreateToken(userID, sessionID)

		assert.Error(t, err) // Verifica se ocorreu o erro esperado
		assert.Empty(t, token)
	})
}

func TestTokenService_ExtractSessionID(t *testing.T) {
	// Configurar chave privada e pública para o teste
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	config.Env.PrivateKey = "privateKey"
	config.Env.PublicKey = "PublicKey"
	config.Env.Cache.SessionExp = 1 // 1 hora para expiração

	tokenService := &tokenService{}

	t.Run("should extract session ID successfully from a valid token", func(t *testing.T) {
		userID := uuid.New()
		sessionID := uuid.New()

		claims := jwt.MapClaims{
			"userId": userID.String(),
			"exp":    time.Now().Add(time.Hour).Unix(),
			"sid":    sessionID.String(),
			"iat":    time.Now().Unix(),
			"iss":    "level-up.com",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
		signedToken, err := token.SignedString(privateKey)

		assert.NoError(t, err)

		extractedSessionID, err := tokenService.ExtractSessionID(signedToken)

		assert.NoError(t, err)
		assert.Equal(t, sessionID, extractedSessionID)
	})

	t.Run("should return an error for invalid token signature", func(t *testing.T) {
		userID := uuid.New()
		sessionID := uuid.New()

		claims := jwt.MapClaims{
			"userId": userID.String(),
			"exp":    time.Now().Add(time.Hour).Unix(),
			"sid":    sessionID.String(),
			"iat":    time.Now().Unix(),
			"iss":    "level-up.com",
		}

		// Criar token com chave privada diferente
		otherPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		assert.NoError(t, err)

		token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
		signedToken, err := token.SignedString(otherPrivateKey)

		assert.NoError(t, err)

		extractedSessionID, err := tokenService.ExtractSessionID(signedToken)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, extractedSessionID)
	})

	t.Run("should return an error for invalid session ID", func(t *testing.T) {
		claims := jwt.MapClaims{
			"userId": uuid.New().String(),
			"exp":    time.Now().Add(time.Hour).Unix(),
			"sid":    "invalid-session-id", // Session ID inválido
			"iat":    time.Now().Unix(),
			"iss":    "level-up.com",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
		signedToken, err := token.SignedString(privateKey)

		assert.NoError(t, err)

		extractedSessionID, err := tokenService.ExtractSessionID(signedToken)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, extractedSessionID)
	})

	t.Run("should return an error for expired token", func(t *testing.T) {
		userID := uuid.New()
		sessionID := uuid.New()

		claims := jwt.MapClaims{
			"userId": userID.String(),
			"exp":    time.Now().Add(-time.Hour).Unix(), // Token expirado
			"sid":    sessionID.String(),
			"iat":    time.Now().Unix(),
			"iss":    "level-up.com",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
		signedToken, err := token.SignedString(privateKey)

		assert.NoError(t, err)

		extractedSessionID, err := tokenService.ExtractSessionID(signedToken)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, extractedSessionID)
	})

	t.Run("should return an error for malformed token", func(t *testing.T) {
		malformedToken := "this-is-not-a-valid-token"

		extractedSessionID, err := tokenService.ExtractSessionID(malformedToken)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, extractedSessionID)
	})
}
