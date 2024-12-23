package models

import "errors"

var (
	ErrMagicLinkNotFound = errors.New("magic link not found")
)

type SignInPayload struct {
	Email string `json:"email" validate:"required,email,max=255"`
}
