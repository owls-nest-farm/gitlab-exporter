package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xanzy/go-gitlab"
)

type Webhook struct {
	PayloadURL            string
	ContentType           string
	EventTypes            []string
	EnableSSLVerification bool
	Active                bool
}

func getWebhooks(project *gitlab.Project) ([]Webhook, error) {
	api := fmt.Sprintf("projects/%d/hooks", project.ID)
	resp, err := newRequest(api)
	if err != nil {
		log.Fatalf("Failed to create new request: %v", err)
	}
	var w []Webhook
	err = json.NewDecoder(resp.Body).Decode(&w)
	return w, err
}
