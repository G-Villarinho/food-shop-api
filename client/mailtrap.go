package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/G-Villarinho/food-shop-api/config"
	"github.com/G-Villarinho/food-shop-api/internal"
)

type MailtrapRecipient struct {
	Email string `json:"email"`
}

type MailtrapPayload struct {
	From struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"from"`
	To       []MailtrapRecipient `json:"to"`
	Subject  string              `json:"subject"`
	Text     string              `json:"text"`
	Html     string              `json:"html"`
	Category string              `json:"category"`
}

type MailtrapResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

//go:generate mockery --name=MailtrapClient --output=../mocks --outpkg=mocks
type MailtrapClient interface {
	SendEmail(ctx context.Context, payload MailtrapPayload) error
}

type mailtrapClient struct {
	di *internal.Di
}

func NewMailtrapClient(di *internal.Di) (MailtrapClient, error) {
	return &mailtrapClient{
		di: di,
	}, nil
}

func (m *mailtrapClient) SendEmail(ctx context.Context, payload MailtrapPayload) error {
	payloadBytes, err := jsoniter.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", config.Env.Email.EmailClientBaseURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Env.Email.EmailClientApiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()

	var genericResponse map[string]interface{}
	if err := jsoniter.NewDecoder(res.Body).Decode(&genericResponse); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("mailtrap API error: %v", genericResponse)
	}

	return nil
}
