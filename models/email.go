package models

type EmailTemplate string

const (
	SignInMagicLink EmailTemplate = "sign-in-magic-link"
)

type Email struct {
	From     string
	FromName string
	To       []string
	Subject  string
	Text     string
	Html     string
}

type EmailResponse struct {
	Status  string
	Message string
}

type EmailQueueTask struct {
	Template EmailTemplate
	Subject  string
	To       []string
	Params   map[string]string
}
