package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type WebhookService struct {
	*BaseService
}

type Webhook struct {
	PayloadURL            string
	ContentType           string
	EventTypes            []string
	EnableSSLVerification bool
	Active                bool
}

func NewWebhookService(e *Exporter) *WebhookService {
	return &WebhookService{
		BaseService: &BaseService{
			exporter: e,
			filename: "webhooks.json",
		},
	}
}

func (ws *WebhookService) GetAll() ([]Webhook, error) {
	api := fmt.Sprintf("projects/%d/hooks", ws.exporter.CurrentProject.ID)
	resp, err := ws.exporter.NewRequest(api)
	if err != nil {
		log.Fatalf("Failed to create new request: %v", err)
	}
	var w []Webhook
	err = json.NewDecoder(resp.Body).Decode(&w)
	return w, err
}
