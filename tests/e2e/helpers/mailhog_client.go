package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Email representeert een email in MailHog
type Email struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// MailhogClient is een helper voor het ophalen van emails uit MailHog
type MailhogClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewMailhogClient maakt een nieuwe MailHog client
func NewMailhogClient(baseURL string) *MailhogClient {
	return &MailhogClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// GetEmailsTo haalt emails op die verzonden zijn naar een specifiek email adres
func (c *MailhogClient) GetEmailsTo(email string) ([]Email, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/v2/search?kind=to&query=%s", c.baseURL, url.QueryEscape(email)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Items []Email `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}
