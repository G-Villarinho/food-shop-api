package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/G-Villarinho/level-up-api/client"
	"github.com/G-Villarinho/level-up-api/config"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/templates"
	jsoniter "github.com/json-iterator/go"
)

type EmailService interface {
	SendEmail(ctx context.Context, task models.EmailQueueTask) error
	SendEmailAsync(ctx context.Context, task models.EmailQueueTask)
}

type emailService struct {
	di              *internal.Di
	emailClient     client.MailtrapClient
	queueService    QueueService
	templateService templates.TemplateService
}

func NewEmailService(di *internal.Di) (EmailService, error) {
	emailClient, err := internal.Invoke[client.MailtrapClient](di)
	if err != nil {
		return nil, err
	}

	queueService, err := internal.Invoke[QueueService](di)
	if err != nil {
		return nil, err
	}

	templateService, err := internal.Invoke[templates.TemplateService](di)
	if err != nil {
		return nil, err
	}

	return &emailService{
		di:              di,
		emailClient:     emailClient,
		queueService:    queueService,
		templateService: templateService,
	}, nil
}

func (e *emailService) SendEmail(ctx context.Context, task models.EmailQueueTask) error {
	content, err := e.templateService.RenderTemplate(string(task.Template), task.Params)
	if err != nil {
		return fmt.Errorf("render %s.html email template: %w", task.Template, err)
	}

	email := models.Email{
		From:     config.Env.Email.EmailSender,
		FromName: "level up auth",
		To:       task.To,
		Subject:  task.Subject,
		Html:     content,
	}

	if err := e.emailClient.SendEmail(ctx, toMailtrapPayload(email)); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func (e *emailService) SendEmailAsync(ctx context.Context, task models.EmailQueueTask) {
	log := slog.With(
		slog.String("service", "email"),
		slog.String("func", "SendEmailAsync"),
	)

	go func() {
		message, err := jsoniter.Marshal(task)
		if err != nil {
			log.Error("marshal email task", slog.String("error", err.Error()))
			return
		}

		if err := e.queueService.Publish(QueueSendEmail, message); err != nil {
			log.Error("publish email task", slog.String("error", err.Error()))
			return
		}
	}()
}

func toMailtrapPayload(email models.Email) client.MailtrapPayload {
	toRecipients := make([]client.MailtrapRecipient, len(email.To))
	for i, recipient := range email.To {
		toRecipients[i] = client.MailtrapRecipient{Email: recipient}
	}

	payload := client.MailtrapPayload{
		To:       toRecipients,
		Subject:  email.Subject,
		Text:     fmt.Sprintf("Plain text fallback for %s", email.Subject),
		Html:     email.Html,
		Category: "Transactional",
	}
	payload.From.Email = email.From
	payload.From.Name = email.FromName

	return payload
}
