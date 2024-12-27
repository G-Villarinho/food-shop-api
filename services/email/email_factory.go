package email

import (
	"github.com/G-Villarinho/food-shop-api/models"
)

type EmailFactory struct{}

func NewEmailTaskFactory() *EmailFactory {
	return &EmailFactory{}
}

func (f *EmailFactory) CreateSignInMagicLinkEmail(to string, name string, magicLink string) models.EmailQueueTask {
	return models.EmailQueueTask{
		To:       []string{to},
		Subject:  "Sign in to Level Up",
		Template: models.SignInMagicLink,
		Params: map[string]string{
			"magic_link": magicLink,
			"name":       name,
		},
	}
}
