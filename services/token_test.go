package services

import (
	"testing"
	"time"

	"github.com/G-Villarinho/food-shop-api/config"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTokenService_CreateToken(t *testing.T) {
	privateKey := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIlmFjHouHxEJrhWvTzdSDsMegE0sEDjd91Dy7DsDEUqoAoGCCqGSM49
AwEHoUQDQgAEmSIeG8APOh4v1rX6ERinckdljyHErmfWPkwpNJOhOMD+8vKAfIbl
dQpPG+XUKikw5UnzVTYP+uKN6kK3XHOTYg==
-----END EC PRIVATE KEY-----`

	config.Env.PrivateKey = privateKey
	config.Env.Cache.SessionExp = 6 // hour

	tokenService := &tokenService{
		di: internal.NewDi(),
	}

	t.Run("should create a valid token successfully", func(t *testing.T) {
		userID := uuid.New()
		sessionID := uuid.New()

		token, err := tokenService.CreateToken(userID, sessionID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, claims := validateToken(token, config.Env.PrivateKey, t)

		assert.True(t, parsedToken.Valid)
		assert.Equal(t, userID.String(), claims["userId"])
		assert.Equal(t, sessionID.String(), claims["sid"])
		assert.Equal(t, "level-up.com", claims["iss"])
	})

	t.Run("should return error when private key is invalid", func(t *testing.T) {
		invalidPrivateKey := `-----BEGIN EC PRIVATE KEY-----
INVALID_KEY
-----END EC PRIVATE KEY-----`
		config.Env.PrivateKey = invalidPrivateKey

		userID := uuid.New()
		sessionID := uuid.New()

		token, err := tokenService.CreateToken(userID, sessionID)
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestTokenService_ExtractSessionID(t *testing.T) {
	privateKey := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIlmFjHouHxEJrhWvTzdSDsMegE0sEDjd91Dy7DsDEUqoAoGCCqGSM49
AwEHoUQDQgAEmSIeG8APOh4v1rX6ERinckdljyHErmfWPkwpNJOhOMD+8vKAfIbl
dQpPG+XUKikw5UnzVTYP+uKN6kK3XHOTYg==
-----END EC PRIVATE KEY-----`

	publicKey := `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEmSIeG8APOh4v1rX6ERinckdljyHE
rmfWPkwpNJOhOMD+8vKAfIbldQpPG+XUKikw5UnzVTYP+uKN6kK3XHOTYg==
-----END PUBLIC KEY-----`

	config.Env.PrivateKey = privateKey
	config.Env.PublicKey = publicKey

	tokenService := &tokenService{
		di: internal.NewDi(),
	}

	t.Run("should extract session ID successfully from a valid token", func(t *testing.T) {
		userID := uuid.New()
		sessionID := uuid.New()
		config.Env.Cache.SessionExp = 6

		token, err := tokenService.CreateToken(userID, sessionID)
		assert.NoError(t, err)

		extractedSessionID, err := tokenService.ExtractSessionID(token)
		assert.NoError(t, err)
		assert.Equal(t, sessionID, extractedSessionID)
	})

	t.Run("should return error when token is invalid", func(t *testing.T) {
		invalidToken := "INVALID_TOKEN"

		sessionID, err := tokenService.ExtractSessionID(invalidToken)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, sessionID)
	})

	t.Run("should return error when public key is invalid", func(t *testing.T) {
		config.Env.PublicKey = `-----BEGIN PUBLIC KEY-----
INVALID_KEY
-----END PUBLIC KEY-----`

		userID := uuid.New()
		sessionID := uuid.New()

		token, err := tokenService.CreateToken(userID, sessionID)
		assert.NoError(t, err)

		sessionID, err = tokenService.ExtractSessionID(token)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, sessionID)
	})

	t.Run("should return error when session ID is missing in token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"userId": uuid.New().String(),
			"exp":    time.Now().Add(time.Hour).Unix(),
			"iat":    time.Now().Unix(),
			"iss":    "level-up.com",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
		privateKey, err := parseECPrivateKey(config.Env.PrivateKey)
		assert.NoError(t, err)

		signedToken, err := token.SignedString(privateKey)
		assert.NoError(t, err)

		sessionID, err := tokenService.ExtractSessionID(signedToken)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, sessionID)
	})
}

func validateToken(token, privateKey string, t *testing.T) (*jwt.Token, jwt.MapClaims) {
	parsedKey, err := parseECPrivateKey(privateKey)
	assert.NoError(t, err)

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return parsedKey.Public(), nil
	})

	assert.NoError(t, err)
	return parsedToken, claims
}
