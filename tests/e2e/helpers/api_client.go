package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// APIClient is een helper voor het maken van HTTP requests naar de API
type APIClient struct {
	BaseURL    string
	httpClient *http.Client
	Headers    map[string]string
}

// NewAPIClient maakt een nieuwe API client
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL:    baseURL,
		httpClient: &http.Client{},
		Headers:    make(map[string]string),
	}
}

// NewAuthenticatedAPIClient maakt een nieuwe geauthenticeerde API client
func NewAuthenticatedAPIClient(baseURL, email, password string) *APIClient {
	client := NewAPIClient(baseURL)

	// Login en token ophalen
	loginReq := map[string]string{
		"email":    email,
		"password": password,
	}

	resp, err := client.Post("/api/auth/login", loginReq)
	if err != nil {
		panic(fmt.Sprintf("Failed to login: %v", err))
	}
	defer resp.Body.Close()

	var loginResp struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		panic(fmt.Sprintf("Failed to decode login response: %v", err))
	}

	client.Headers["Authorization"] = "Bearer " + loginResp.AccessToken
	return client
}

// Post doet een POST request naar de API
func (c *APIClient) Post(path string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	return c.httpClient.Do(req)
}

// Get doet een GET request naar de API
func (c *APIClient) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	return c.httpClient.Do(req)
}
